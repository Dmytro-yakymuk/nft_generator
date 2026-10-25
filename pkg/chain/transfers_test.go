package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeNode answers the JSON-RPC requests the reader sends, single and batched, and records them.
type fakeNode struct {
	mu       sync.Mutex
	logs     []map[string]interface{}
	txValues map[common.Hash]*big.Int
	ranges   [][2]uint64
	methods  []string
	// head, if set, is the latest block of a node behind the block number other requests saw.
	head uint64
}

type rpcRequest struct {
	ID     json.RawMessage   `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func (node *fakeNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")

	if strings.HasPrefix(strings.TrimSpace(string(body)), "[") {
		var requests []rpcRequest
		_ = json.Unmarshal(body, &requests)
		responses := make([]map[string]interface{}, len(requests))
		for i, request := range requests {
			responses[i] = node.answer(request)
		}
		_ = json.NewEncoder(w).Encode(responses)
		return
	}

	var request rpcRequest
	_ = json.Unmarshal(body, &request)
	_ = json.NewEncoder(w).Encode(node.answer(request))
}

func (node *fakeNode) answer(request rpcRequest) map[string]interface{} {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.methods = append(node.methods, request.Method)

	var result interface{}
	switch request.Method {
	case "eth_blockNumber":
		result = hexutil.Uint64(130)
	case "eth_getLogs":
		var filter struct {
			FromBlock hexutil.Uint64 `json:"fromBlock"`
			ToBlock   hexutil.Uint64 `json:"toBlock"`
		}
		_ = json.Unmarshal(request.Params[0], &filter)
		node.ranges = append(node.ranges, [2]uint64{uint64(filter.FromBlock), uint64(filter.ToBlock)})
		if node.head != 0 && uint64(filter.ToBlock) > node.head {
			return map[string]interface{}{"jsonrpc": "2.0", "id": request.ID, "error": map[string]interface{}{
				"code":    -32000,
				"message": fmt.Sprintf("block range extends beyond current head block: requested %d, head %d", filter.ToBlock, node.head),
			}}
		}

		logs := []map[string]interface{}{}
		for _, log := range node.logs {
			block, _ := hexutil.DecodeUint64(log["blockNumber"].(string))
			if block >= uint64(filter.FromBlock) && block <= uint64(filter.ToBlock) {
				logs = append(logs, log)
			}
		}
		result = logs
	case "eth_getBlockByNumber":
		var number hexutil.Uint64
		_ = json.Unmarshal(request.Params[0], &number)
		result = map[string]interface{}{
			"hash":      common.BigToHash(big.NewInt(int64(number))),
			"timestamp": hexutil.Uint64(1760000000 + uint64(number)*12),
		}
	case "eth_getTransactionByHash":
		var hash common.Hash
		_ = json.Unmarshal(request.Params[0], &hash)
		if value, ok := node.txValues[hash]; ok {
			result = map[string]interface{}{"hash": hash, "value": (*hexutil.Big)(value)}
		}
	}

	return map[string]interface{}{"jsonrpc": "2.0", "id": request.ID, "result": result}
}

func (node *fakeNode) count(method string) int {
	node.mu.Lock()
	defer node.mu.Unlock()
	count := 0
	for _, m := range node.methods {
		if m == method {
			count++
		}
	}
	return count
}

// transferLog is an ERC-721 Transfer log; timestamp 0 means the node does not return blockTimestamp.
func transferLog(block uint64, timestamp uint64, tx common.Hash, from, to common.Address, tokenID int64) map[string]interface{} {
	log := map[string]interface{}{
		"topics": []common.Hash{
			TransferTopic,
			common.BytesToHash(from.Bytes()),
			common.BytesToHash(to.Bytes()),
			common.BigToHash(big.NewInt(tokenID)),
		},
		"blockNumber":     hexutil.Uint64(block),
		"transactionHash": tx,
		"logIndex":        hexutil.Uint(0),
		"removed":         false,
	}
	if timestamp != 0 {
		log["blockTimestamp"] = hexutil.Uint64(timestamp)
	}
	encoded, _ := json.Marshal(log)
	var decoded map[string]interface{}
	_ = json.Unmarshal(encoded, &decoded)
	return decoded
}

func TestReader(t *testing.T) {
	ctx := context.Background()
	contract := common.HexToAddress("0x2df71a8114129312c0e5e26c26a9bc61e22d6fd7")
	buyer := common.HexToAddress("0xFf11E1F8721dadd907e8A9519B8111E20717c580")
	mintTx, giftTx := common.HexToHash("0x01"), common.HexToHash("0x02")
	price, _ := new(big.Int).SetString("2000000000000000", 10)

	erc20 := transferLog(105, 1760000100, giftTx, buyer, contract, 0)
	erc20["topics"] = erc20["topics"].([]interface{})[:3]
	removed := transferLog(105, 1760000100, giftTx, buyer, contract, 9)
	removed["removed"] = true

	node := &fakeNode{
		logs: []map[string]interface{}{
			transferLog(101, 1760000000, mintTx, common.Address{}, buyer, 0),
			transferLog(101, 1760000000, mintTx, common.Address{}, buyer, 1),
			erc20,
			removed,
			// this node does not return blockTimestamp for the block 125.
			transferLog(125, 0, giftTx, buyer, contract, 1),
		},
		txValues: map[common.Hash]*big.Int{mintTx: price, giftTx: big.NewInt(0)},
	}
	server := httptest.NewServer(node)
	defer server.Close()

	reader, err := NewReader(server.URL, 10)
	require.NoError(t, err)
	defer reader.Close()

	transfers, err := reader.Transfers(ctx, contract, 100, 125)
	require.NoError(t, err)

	// the range is read in chunks of 10 blocks.
	assert.Equal(t, [][2]uint64{{100, 109}, {110, 119}, {120, 125}}, node.ranges)

	require.Len(t, transfers, 3)
	assert.Equal(t, common.Address{}, transfers[0].From)
	assert.Equal(t, buyer, transfers[0].To)
	assert.Equal(t, int64(1), transfers[1].TokenID.Int64())
	assert.Equal(t, time.Unix(1760000000, 0).UTC(), transfers[0].Timestamp)
	assert.Equal(t, price, transfers[0].TxValue)
	assert.Equal(t, 2, transfers[0].TxTransfers)

	assert.Equal(t, contract, transfers[2].To)
	assert.Equal(t, time.Unix(1760000000+125*12, 0).UTC(), transfers[2].Timestamp)
	assert.Zero(t, transfers[2].TxValue.Sign())
	assert.Equal(t, 1, transfers[2].TxTransfers)
	assert.Equal(t, 1, node.count("eth_getBlockByNumber"))
	assert.Equal(t, 2, node.count("eth_getTransactionByHash"))

	t.Run("transaction values are cached", func(t *testing.T) {
		_, err := reader.Transfers(ctx, contract, 100, 109)
		require.NoError(t, err)
		assert.Equal(t, 2, node.count("eth_getTransactionByHash"))
	})

	t.Run("unknown transaction", func(t *testing.T) {
		node.logs = append(node.logs, transferLog(130, 1760000200, common.HexToHash("0x03"), buyer, contract, 0))
		_, err := reader.Transfers(ctx, contract, 130, 130)
		require.Error(t, err)
		assert.True(t, ErrChain.Has(err))
	})

	t.Run("a node behind the requested range", func(t *testing.T) {
		node.mu.Lock()
		node.head, node.ranges = 112, nil
		node.mu.Unlock()
		defer func() { node.head = 0 }()

		// the range stops at the head of the node: the blocks after it are read on the next refresh.
		transfers, err := reader.Transfers(ctx, contract, 100, 125)
		require.NoError(t, err)
		assert.Equal(t, [][2]uint64{{100, 109}, {110, 119}, {110, 112}}, node.ranges)
		assert.Len(t, transfers, 2)

		node.mu.Lock()
		node.ranges = nil
		node.mu.Unlock()
		transfers, err = reader.Transfers(ctx, contract, 120, 125)
		require.NoError(t, err)
		assert.Equal(t, [][2]uint64{{120, 125}}, node.ranges)
		assert.Empty(t, transfers)
	})

	t.Run("block number", func(t *testing.T) {
		number, err := reader.BlockNumber(ctx)
		require.NoError(t, err)
		assert.Equal(t, uint64(130), number)
	})
}
