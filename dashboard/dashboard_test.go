package dashboard

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/collections"
	"generator/internal/logger/zaplog"
	"generator/internal/nft"
	"generator/pkg/chain"
)

var (
	contract = common.HexToAddress("0x2df71a8114129312c0e5e26c26a9bc61e22d6fd7")
	owner    = common.HexToAddress("0xa7aF8410b61CA79583F9995963fa947DA52cce4e")
	buyer    = common.HexToAddress("0xFf11E1F8721dadd907e8A9519B8111E20717c580")
	friend   = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	start    = time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)
)

// wei returns the amount of ether in wei.
func wei(ether float64) *big.Int {
	value, _ := new(big.Float).Mul(big.NewFloat(ether), weiInEther).Int(nil)
	return value
}

// transfer returns a transfer of the token in the block; a block is minted every minute.
func transfer(block uint64, from, to common.Address, tokenID int64, txValue *big.Int, txTransfers int) chain.Transfer {
	return chain.Transfer{
		From:        from,
		To:          to,
		TokenID:     big.NewInt(tokenID),
		BlockNumber: block,
		Timestamp:   start.Add(time.Duration(block) * time.Minute),
		TxHash:      common.BigToHash(big.NewInt(int64(block))),
		TxValue:     txValue,
		TxTransfers: txTransfers,
	}
}

// history is the activity of the demo collection: three paid mints, a free mint,
// a resale, a gift and a burn.
func history() []chain.Transfer {
	zero := common.Address{}
	return []chain.Transfer{
		transfer(10, zero, owner, 0, wei(0.003), 3),
		transfer(10, zero, owner, 1, wei(0.003), 3),
		transfer(10, zero, owner, 2, wei(0.003), 3),
		transfer(20, zero, buyer, 3, big.NewInt(0), 1),
		transfer(30, owner, buyer, 0, wei(0.01), 1),
		transfer(40, buyer, friend, 3, big.NewInt(0), 1),
		transfer(50, owner, zero, 2, big.NewInt(0), 1),
	}
}

func TestToEvents(t *testing.T) {
	events := toEvents(history())
	require.Len(t, events, 7)

	types := make([]string, len(events))
	for i, event := range events {
		types[i] = event.Type
		assert.Equal(t, Currency, event.Currency)
	}
	assert.Equal(t, []string{EventMint, EventMint, EventMint, EventMint, EventSale, EventTransfer, EventBurn}, types)

	// mint(3) for 0.003 ETH: each token costs 0.001 ETH.
	assert.Equal(t, "0.001", events[0].Price.Text('g', 10))
	assert.Equal(t, "0", events[3].Price.Text('g', 10))
	assert.Equal(t, "0.01", events[4].Price.Text('g', 10))
	assert.Zero(t, events[6].Price.Sign())
}

func TestCollectionStats(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, CollectionStats{}, collectionStats(nil))
	})

	t.Run("history", func(t *testing.T) {
		// token 2 is burned; tokens 0, 1, 3 belong to buyer, owner and friend; 4 paid events.
		assert.Equal(t, CollectionStats{Items: 3, Owners: 3, Sales: 4}, collectionStats(history()))
	})
}

func TestSalesSummary(t *testing.T) {
	events := toEvents(history())

	t.Run("all time", func(t *testing.T) {
		summary := salesSummary(events, time.Time{}, time.Time{})
		require.Len(t, summary.Events, 4)
		assert.Equal(t, "0.001", summary.FloorPrice.Text('g', 10))
		assert.Equal(t, "0.01", summary.HighestPrice.Text('g', 10))
		assert.Equal(t, "0.013", summary.Volume.Text('g', 10))
		assert.Equal(t, "0.00325", summary.AveragePrice.Text('g', 10))
	})

	t.Run("period", func(t *testing.T) {
		summary := salesSummary(events, start.Add(25*time.Minute), start.Add(35*time.Minute))
		require.Len(t, summary.Events, 1)
		assert.Equal(t, "0.01", summary.FloorPrice.Text('g', 10))
		assert.Equal(t, "0.01", summary.AveragePrice.Text('g', 10))
	})

	t.Run("no sales", func(t *testing.T) {
		summary := salesSummary(events, start.Add(time.Hour), time.Time{})
		assert.Empty(t, summary.Events)
		for _, price := range []*big.Float{summary.FloorPrice, summary.HighestPrice, summary.AveragePrice, summary.Volume} {
			require.NotNil(t, price)
			assert.Zero(t, price.Sign())
		}
	})
}

// fakeChain is a node with the given transfers that counts requests.
type fakeChain struct {
	latest    uint64
	transfers []chain.Transfer
	requests  [][2]uint64
}

func (fake *fakeChain) BlockNumber(ctx context.Context) (uint64, error) {
	return fake.latest, nil
}

func (fake *fakeChain) Transfers(ctx context.Context, contract common.Address, from, to uint64) ([]chain.Transfer, error) {
	fake.requests = append(fake.requests, [2]uint64{from, to})
	var result []chain.Transfer
	for _, transfer := range fake.transfers {
		if transfer.BlockNumber >= from && transfer.BlockNumber <= to {
			result = append(result, transfer)
		}
	}
	return result, nil
}

// fakeCollections stores collections in memory; token images are "image-<tokenId>".
type fakeCollections map[uuid.UUID]collections.Collection

func (fake fakeCollections) Get(ctx context.Context, id uuid.UUID) (collections.Collection, error) {
	project, ok := fake[id]
	if !ok {
		return project, collections.ErrNoCollection.New("")
	}
	return project, nil
}

func (fake fakeCollections) TokenMetadata(ctx context.Context, id uuid.UUID, tokenID int) (nft.Metadata, error) {
	return nft.Metadata{Image: fmt.Sprintf("image-%d", tokenID)}, nil
}

func TestService(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	deployed := collections.Collection{ID: uuid.New(), UserID: userID, ContractAddress: contract, ContractBlock: 5}
	notDeployed := collections.Collection{ID: uuid.New(), UserID: userID}
	node := &fakeChain{latest: 60, transfers: history()}
	service := NewService(Config{}, node, fakeCollections{deployed.ID: deployed, notDeployed.ID: notDeployed}, zaplog.NewLog())
	now := start.Add(time.Hour)
	service.now = func() time.Time { return now }

	t.Run("access", func(t *testing.T) {
		_, err := service.CollectionStats(ctx, uuid.New(), userID)
		assert.True(t, collections.ErrNoCollection.Has(err))

		_, err = service.CollectionStats(ctx, deployed.ID, uuid.New())
		assert.True(t, collections.ErrCollectionOwnership.Has(err))

		_, err = service.Events(ctx, notDeployed.ID, userID, 1)
		assert.True(t, collections.ErrContractNotDeployed.Has(err))

		assert.Empty(t, node.requests)
	})

	t.Run("stats are read from the deployment block", func(t *testing.T) {
		stats, err := service.CollectionStats(ctx, deployed.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, CollectionStats{Items: 3, Owners: 3, Sales: 4}, stats)
		assert.Equal(t, [][2]uint64{{5, 60}}, node.requests)
	})

	t.Run("events newest first", func(t *testing.T) {
		page, err := service.Events(ctx, deployed.ID, userID, 1)
		require.NoError(t, err)
		require.Len(t, page.Events, 7)
		assert.Equal(t, EventBurn, page.Events[0].Type)
		assert.Equal(t, EventMint, page.Events[6].Type)
		assert.Equal(t, "image-2", page.Events[0].Image)
		assert.Empty(t, page.NextPage)
		assert.Empty(t, page.PreviousPage)

		_, err = service.Events(ctx, deployed.ID, userID, 0)
		require.Error(t, err)

		// three requests within the refresh interval read the node once.
		assert.Len(t, node.requests, 1)
	})

	t.Run("sales of the last day", func(t *testing.T) {
		summary, err := service.EventDataAndAveragePrices(ctx, deployed.ID, userID, 1, 0, 0)
		require.NoError(t, err)
		assert.Len(t, summary.Events, 4)
		assert.Equal(t, "0.013", summary.Volume.Text('g', 10))

		summary, err = service.EventDataAndAveragePrices(ctx, deployed.ID, userID, 0, start.Add(25*time.Minute).Unix(), 0)
		require.NoError(t, err)
		assert.Len(t, summary.Events, 1)
	})

	t.Run("refresh reads new blocks and replaces reorganized ones", func(t *testing.T) {
		// the block 50 with the burn was replaced: the token was given to a friend instead.
		transfers := history()[:6]
		transfers = append(transfers,
			transfer(50, owner, friend, 2, big.NewInt(0), 1),
			transfer(70, common.Address{}, buyer, 4, wei(0.002), 2),
			transfer(70, common.Address{}, buyer, 5, wei(0.002), 2),
		)
		node.transfers = transfers
		node.latest = 75
		now = now.Add(refreshInterval)

		stats, err := service.CollectionStats(ctx, deployed.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, CollectionStats{Items: 6, Owners: 3, Sales: 6}, stats)
		assert.Equal(t, [2]uint64{60 + 1 - reorgDepth, 75}, node.requests[len(node.requests)-1])

		// a node behind the one read last time does not reset the history.
		node.latest = 70
		now = now.Add(refreshInterval)
		stats, err = service.CollectionStats(ctx, deployed.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, 6, stats.Items)
		assert.Len(t, node.requests, 2)
	})

	t.Run("pages", func(t *testing.T) {
		node.transfers = nil
		for i := 0; i < 2*PageSize+3; i++ {
			node.transfers = append(node.transfers, transfer(uint64(100+i), common.Address{}, buyer, int64(i), wei(0.001), 1))
		}
		other := collections.Collection{ID: uuid.New(), UserID: userID, ContractAddress: common.HexToAddress("0x01"), ContractBlock: 100}
		service.collections = fakeCollections{other.ID: other}
		node.latest = 200

		first, err := service.Events(ctx, other.ID, userID, 1)
		require.NoError(t, err)
		require.Len(t, first.Events, PageSize)
		assert.Equal(t, strconv.Itoa(2*PageSize+2), first.Events[0].TokenID)
		assert.Equal(t, "2", first.NextPage)
		assert.Empty(t, first.PreviousPage)

		last, err := service.Events(ctx, other.ID, userID, 3)
		require.NoError(t, err)
		require.Len(t, last.Events, 3)
		assert.Equal(t, "0", last.Events[2].TokenID)
		assert.Empty(t, last.NextPage)
		assert.Equal(t, "2", last.PreviousPage)

		beyond, err := service.Events(ctx, other.ID, userID, 4)
		require.NoError(t, err)
		assert.Empty(t, beyond.Events)
		assert.Equal(t, "3", beyond.PreviousPage)
	})
}
