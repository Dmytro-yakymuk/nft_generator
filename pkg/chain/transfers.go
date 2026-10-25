package chain

import (
	"context"
	"math/big"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zeebo/errs"
)

// ErrChain indicates that the node could not return the requested data.
var ErrChain = errs.Class("chain")

// TransferTopic is the signature hash of the ERC-721 Transfer(address,address,uint256) event.
var TransferTopic = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

// DefaultLogsBlockRange is the widest block range of one eth_getLogs request
// that public nodes accept (publicnode.com rejects more than 50000 blocks).
const DefaultLogsBlockRange = 10000

// Transfer is an ERC-721 Transfer event together with the ether paid in its transaction.
type Transfer struct {
	From        common.Address
	To          common.Address
	TokenID     *big.Int
	BlockNumber uint64
	LogIndex    uint
	Timestamp   time.Time
	TxHash      common.Hash
	// TxValue is the ether sent with the transaction, wei.
	TxValue *big.Int
	// TxTransfers is the number of the contract transfers in the transaction:
	// TxValue / TxTransfers is the price of one token when several are minted at once.
	TxTransfers int
}

// Reader reads ERC-721 events through the JSON-RPC API of a node.
type Reader struct {
	client     *rpc.Client
	blockRange uint64

	// transactions never change once mined, so their values are kept for the next requests.
	mu       sync.Mutex
	txValues map[common.Hash]*big.Int
}

// NewReader returns a reader of the node at rpcURL; the connection is opened on the first request.
func NewReader(rpcURL string, blockRange uint64) (*Reader, error) {
	client, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, ErrChain.Wrap(err)
	}
	if blockRange == 0 {
		blockRange = DefaultLogsBlockRange
	}

	return &Reader{client: client, blockRange: blockRange, txValues: make(map[common.Hash]*big.Int)}, nil
}

// Close closes the connection to the node.
func (reader *Reader) Close() {
	reader.client.Close()
}

// BlockNumber returns the number of the latest block.
func (reader *Reader) BlockNumber(ctx context.Context) (uint64, error) {
	var number hexutil.Uint64
	err := reader.client.CallContext(ctx, &number, "eth_blockNumber")
	return uint64(number), ErrChain.Wrap(err)
}

// behindHead matches the error of a node whose latest block is behind the requested range
// (geth and reth: "block range extends beyond current head block: requested 107, head 106").
var behindHead = regexp.MustCompile(`beyond current head block.*head (\d+)`)

// nodeHead returns the latest block of the node if err says the requested range is beyond it.
func nodeHead(err error) (uint64, bool) {
	if err == nil {
		return 0, false
	}
	match := behindHead.FindStringSubmatch(err.Error())
	if match == nil {
		return 0, false
	}
	head, parseErr := strconv.ParseUint(match[1], 10, 64)
	return head, parseErr == nil
}

// transferLogs returns the Transfer logs of the contract in blocks [from, to] in one eth_getLogs request.
func (reader *Reader) transferLogs(ctx context.Context, contract common.Address, from, to uint64) ([]rpcLog, error) {
	var logs []rpcLog
	err := reader.client.CallContext(ctx, &logs, "eth_getLogs", map[string]interface{}{
		"address":   contract,
		"fromBlock": hexutil.EncodeUint64(from),
		"toBlock":   hexutil.EncodeUint64(to),
		"topics":    []interface{}{TransferTopic},
	})
	return logs, err
}

// rpcLog is a log entry of eth_getLogs; blockTimestamp is returned by recent nodes (geth, reth, anvil).
type rpcLog struct {
	Topics         []common.Hash   `json:"topics"`
	BlockNumber    hexutil.Uint64  `json:"blockNumber"`
	BlockTimestamp *hexutil.Uint64 `json:"blockTimestamp"`
	TxHash         common.Hash     `json:"transactionHash"`
	LogIndex       hexutil.Uint    `json:"logIndex"`
	Removed        bool            `json:"removed"`
}

// Transfers returns the Transfer events of the contract in blocks [from, to] in chain order.
func (reader *Reader) Transfers(ctx context.Context, contract common.Address, from, to uint64) ([]Transfer, error) {
	var logs []rpcLog
	for start := from; start <= to; start += reader.blockRange {
		end := start + reader.blockRange - 1
		if end > to {
			end = to
		}

		chunk, err := reader.transferLogs(ctx, contract, start, end)
		if head, behind := nodeHead(err); behind && head < end {
			// a load-balanced public node may answer from a server a block or two behind the one that
			// gave the latest block number: the range stops at its head, and the dashboard reads the
			// last blocks again on its next refresh.
			if head < start {
				break
			}
			chunk, err = reader.transferLogs(ctx, contract, start, head)
			to = head
		}
		if err != nil {
			return nil, ErrChain.New("eth_getLogs %d-%d: %v", start, end, err)
		}
		logs = append(logs, chunk...)
	}

	transfers := make([]Transfer, 0, len(logs))
	txTransfers := make(map[common.Hash]int)
	for _, log := range logs {
		// ERC-20 Transfer has the same signature but 3 topics: the amount is not indexed.
		if log.Removed || len(log.Topics) != 4 {
			continue
		}
		transfer := Transfer{
			From:        common.BytesToAddress(log.Topics[1].Bytes()),
			To:          common.BytesToAddress(log.Topics[2].Bytes()),
			TokenID:     log.Topics[3].Big(),
			BlockNumber: uint64(log.BlockNumber),
			LogIndex:    uint(log.LogIndex),
			TxHash:      log.TxHash,
		}
		if log.BlockTimestamp != nil {
			transfer.Timestamp = time.Unix(int64(*log.BlockTimestamp), 0).UTC()
		}
		transfers = append(transfers, transfer)
		txTransfers[log.TxHash]++
	}

	if err := reader.fillTimestamps(ctx, transfers); err != nil {
		return nil, err
	}

	values, err := reader.transactionValues(ctx, txTransfers)
	if err != nil {
		return nil, err
	}
	for i := range transfers {
		transfers[i].TxValue = values[transfers[i].TxHash]
		transfers[i].TxTransfers = txTransfers[transfers[i].TxHash]
	}

	return transfers, nil
}

// rpcHeader is the part of eth_getBlockByNumber result the reader needs; a zero hash means the node returned null.
type rpcHeader struct {
	Hash      common.Hash    `json:"hash"`
	Timestamp hexutil.Uint64 `json:"timestamp"`
}

// rpcTransaction is the part of eth_getTransactionByHash result the reader needs; a zero hash means the node returned null.
type rpcTransaction struct {
	Hash  common.Hash `json:"hash"`
	Value hexutil.Big `json:"value"`
}

// fillTimestamps reads block times for nodes that do not return blockTimestamp in logs.
func (reader *Reader) fillTimestamps(ctx context.Context, transfers []Transfer) error {
	var batch []rpc.BatchElem
	headers := make(map[uint64]*rpcHeader)
	for _, transfer := range transfers {
		if !transfer.Timestamp.IsZero() || headers[transfer.BlockNumber] != nil {
			continue
		}
		header := new(rpcHeader)
		headers[transfer.BlockNumber] = header
		batch = append(batch, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{hexutil.EncodeUint64(transfer.BlockNumber), false},
			Result: header,
		})
	}
	if err := reader.batchCall(ctx, batch); err != nil {
		return err
	}

	for number, header := range headers {
		if header.Hash == (common.Hash{}) {
			return ErrChain.New("block %d not found", number)
		}
	}
	for i := range transfers {
		if header := headers[transfers[i].BlockNumber]; header != nil {
			transfers[i].Timestamp = time.Unix(int64(header.Timestamp), 0).UTC()
		}
	}

	return nil
}

// transactionValues returns the ether sent with each of the transactions, wei.
func (reader *Reader) transactionValues(ctx context.Context, hashes map[common.Hash]int) (map[common.Hash]*big.Int, error) {
	values := make(map[common.Hash]*big.Int, len(hashes))
	txs := make(map[common.Hash]*rpcTransaction)
	var batch []rpc.BatchElem

	reader.mu.Lock()
	for hash := range hashes {
		if value, ok := reader.txValues[hash]; ok {
			values[hash] = value
			continue
		}
		tx := new(rpcTransaction)
		txs[hash] = tx
		batch = append(batch, rpc.BatchElem{Method: "eth_getTransactionByHash", Args: []interface{}{hash}, Result: tx})
	}
	reader.mu.Unlock()

	if err := reader.batchCall(ctx, batch); err != nil {
		return nil, err
	}

	reader.mu.Lock()
	defer reader.mu.Unlock()
	for hash, tx := range txs {
		if tx.Hash == (common.Hash{}) {
			return nil, ErrChain.New("transaction %s not found", hash.Hex())
		}
		value := tx.Value.ToInt()
		reader.txValues[hash] = value
		values[hash] = value
	}

	return values, nil
}

// batchCall sends the requests in one HTTP round trip and returns the first error.
func (reader *Reader) batchCall(ctx context.Context, batch []rpc.BatchElem) error {
	if len(batch) == 0 {
		return nil
	}
	if err := reader.client.BatchCallContext(ctx, batch); err != nil {
		return ErrChain.Wrap(err)
	}
	for _, elem := range batch {
		if elem.Error != nil {
			return ErrChain.New("%s %v: %v", elem.Method, elem.Args[0], elem.Error)
		}
	}

	return nil
}
