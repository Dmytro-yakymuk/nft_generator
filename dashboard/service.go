package dashboard

import (
	"context"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/collections"
	"generator/internal/logger"
	"generator/internal/nft"
	"generator/pkg/chain"
)

// ErrDashboard indicates that there was an error in the service.
var ErrDashboard = errs.Class("dashboard service error")

// refreshInterval is how long the transfers read from the node are reused:
// the dashboard page sends three requests at once.
const refreshInterval = 10 * time.Second

// reorgDepth is the number of the last read blocks that are read again on refresh,
// so the dashboard drops events of blocks replaced by a chain reorganization
// and picks up events a lagging node had not returned yet.
const reorgDepth = 12

// weiInEther is 10^18.
var weiInEther = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

// Blockchain reads token transfers of collection contracts from the network.
type Blockchain interface {
	// BlockNumber returns the number of the latest block.
	BlockNumber(ctx context.Context) (uint64, error)
	// Transfers returns the Transfer events of the contract in blocks [from, to] in chain order.
	Transfers(ctx context.Context, contract common.Address, from, to uint64) ([]chain.Transfer, error)
}

// Collections provides the collection data the dashboard needs.
type Collections interface {
	// Get returns the collection by id.
	Get(ctx context.Context, id uuid.UUID) (collections.Collection, error)
	// TokenMetadata returns the metadata of the collection token.
	TokenMetadata(ctx context.Context, id uuid.UUID, tokenID int) (nft.Metadata, error)
}

// Service is handling dashboard related logic: it builds collection analytics
// from the Transfer events of the collection contract and the ether paid in their transactions.
//
// architecture: Service
type Service struct {
	config      Config
	chain       Blockchain
	collections Collections
	log         logger.Logger
	now         func() time.Time

	mu       sync.Mutex
	activity map[common.Address]*activity
}

// activity is the transfer history of a contract read from the node.
type activity struct {
	mu        sync.Mutex
	transfers []chain.Transfer
	scanned   bool
	scannedTo uint64
	checkedAt time.Time
}

// NewService is a constructor for dashboard service.
func NewService(config Config, blockchain Blockchain, collections Collections, logger logger.Logger) *Service {
	return &Service{
		config:      config,
		chain:       blockchain,
		collections: collections,
		log:         logger,
		now:         time.Now,
		activity:    make(map[common.Address]*activity),
	}
}

// Events returns a page of the collection events, newest first; pages start from 1.
func (service *Service) Events(ctx context.Context, projectID, userID uuid.UUID, page int) (EventsWithPage, error) {
	if page < 1 {
		return EventsWithPage{}, ErrDashboard.New("page %d", page)
	}

	project, transfers, err := service.collectionTransfers(ctx, projectID, userID)
	if err != nil {
		return EventsWithPage{}, err
	}

	// transfers are in chain order, the table shows the newest first.
	events := toEvents(transfers)
	for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
		events[i], events[j] = events[j], events[i]
	}

	result := EventsWithPage{Events: []Event{}}
	start := (page - 1) * PageSize
	if start < len(events) {
		end := start + PageSize
		if end > len(events) {
			end = len(events)
		}
		result.Events = events[start:end]
		if end < len(events) {
			result.NextPage = strconv.Itoa(page + 1)
		}
	}
	if page > 1 {
		result.PreviousPage = strconv.Itoa(page - 1)
	}

	for i := range result.Events {
		result.Events[i].Image = service.tokenImage(ctx, project.ID, result.Events[i].TokenID)
	}

	return result, nil
}

// CollectionStats returns the number of tokens, owners and sales of the collection.
func (service *Service) CollectionStats(ctx context.Context, projectID, userID uuid.UUID) (CollectionStats, error) {
	_, transfers, err := service.collectionTransfers(ctx, projectID, userID)
	if err != nil {
		return CollectionStats{}, err
	}

	return collectionStats(transfers), nil
}

// EventDataAndAveragePrices returns paid events and their prices for the last days (0 - all time)
// or between the unix times from and to (0 - no limit).
func (service *Service) EventDataAndAveragePrices(ctx context.Context, projectID, userID uuid.UUID, days int, from, to int64) (EventAndAveragePrices, error) {
	_, transfers, err := service.collectionTransfers(ctx, projectID, userID)
	if err != nil {
		return EventAndAveragePrices{}, err
	}

	var after, before time.Time
	if days > 0 {
		after = service.now().AddDate(0, 0, -days)
	}
	if from > 0 {
		after = time.Unix(from, 0)
	}
	if to > 0 {
		before = time.Unix(to, 0)
	}

	return salesSummary(toEvents(transfers), after, before), nil
}

// collectionTransfers returns the transfers of the collection contract after checking
// that the collection belongs to the user and its contract is deployed.
func (service *Service) collectionTransfers(ctx context.Context, projectID, userID uuid.UUID) (collections.Collection, []chain.Transfer, error) {
	project, err := service.collections.Get(ctx, projectID)
	if err != nil {
		return project, nil, ErrDashboard.Wrap(err)
	}
	if project.UserID != userID {
		return project, nil, collections.ErrCollectionOwnership.New("")
	}
	if project.ContractAddress == (common.Address{}) {
		return project, nil, collections.ErrContractNotDeployed.New("")
	}

	transfers, err := service.transfers(ctx, project.ContractAddress, project.ContractBlock)
	return project, transfers, err
}

// transfers returns the transfers of the contract deployed in the block deployBlock. The first call
// reads them from the deployment block; later calls read only new blocks.
func (service *Service) transfers(ctx context.Context, contract common.Address, deployBlock uint64) ([]chain.Transfer, error) {
	service.mu.Lock()
	history, ok := service.activity[contract]
	if !ok {
		history = new(activity)
		service.activity[contract] = history
	}
	service.mu.Unlock()

	history.mu.Lock()
	defer history.mu.Unlock()

	now := service.now()
	if history.scanned && now.Sub(history.checkedAt) < refreshInterval {
		return history.transfers, nil
	}

	latest, err := service.chain.BlockNumber(ctx)
	if err != nil {
		return nil, ErrDashboard.Wrap(err)
	}

	from := deployBlock
	if history.scanned {
		// a node behind the one read last time has nothing new.
		if latest <= history.scannedTo {
			history.checkedAt = now
			return history.transfers, nil
		}
		if history.scannedTo >= deployBlock+reorgDepth {
			from = history.scannedTo + 1 - reorgDepth
		}
	}

	fresh, err := service.chain.Transfers(ctx, contract, from, latest)
	if err != nil {
		return nil, ErrDashboard.Wrap(err)
	}

	// a new slice: callers may still read the previous one.
	transfers := make([]chain.Transfer, 0, len(history.transfers)+len(fresh))
	for _, transfer := range history.transfers {
		if transfer.BlockNumber < from {
			transfers = append(transfers, transfer)
		}
	}
	transfers = append(transfers, fresh...)

	history.transfers, history.scanned, history.scannedTo, history.checkedAt = transfers, true, latest, now

	return transfers, nil
}

// tokenImage returns the image of the token from the metadata the platform serves, or "" if there is none.
func (service *Service) tokenImage(ctx context.Context, projectID uuid.UUID, tokenID string) string {
	id, err := strconv.Atoi(tokenID)
	if err != nil {
		return ""
	}
	metadata, err := service.collections.TokenMetadata(ctx, projectID, id)
	if err != nil {
		return ""
	}

	return metadata.Image
}

// toEvents converts transfers in chain order to events. When several tokens are transferred
// in one paid transaction (e.g. mint(3)), each gets an equal share of the payment.
func toEvents(transfers []chain.Transfer) []Event {
	events := make([]Event, 0, len(transfers))
	for _, transfer := range transfers {
		price := new(big.Float)
		if transfer.TxValue != nil && transfer.TxValue.Sign() > 0 && transfer.TxTransfers > 0 {
			price.Quo(new(big.Float).SetInt(transfer.TxValue), weiInEther)
			price.Quo(price, big.NewFloat(float64(transfer.TxTransfers)))
		}

		eventType := EventTransfer
		switch {
		case transfer.From == (common.Address{}):
			eventType = EventMint
		case transfer.To == (common.Address{}):
			eventType = EventBurn
		case price.Sign() > 0:
			eventType = EventSale
		}

		events = append(events, Event{
			From:        transfer.From,
			To:          transfer.To,
			TokenID:     transfer.TokenID.String(),
			Timestamp:   transfer.Timestamp,
			Type:        eventType,
			Price:       price,
			Currency:    Currency,
			TxHash:      transfer.TxHash,
			BlockNumber: transfer.BlockNumber,
		})
	}

	return events
}

// collectionStats replays transfers in chain order to find the current owner of each token.
func collectionStats(transfers []chain.Transfer) CollectionStats {
	owners := make(map[string]common.Address)
	var stats CollectionStats
	for _, event := range toEvents(transfers) {
		if event.Type == EventBurn {
			delete(owners, event.TokenID)
		} else {
			owners[event.TokenID] = event.To
		}
		if event.Price.Sign() > 0 {
			stats.Sales++
		}
	}

	holders := make(map[common.Address]struct{})
	for _, owner := range owners {
		holders[owner] = struct{}{}
	}
	stats.Items, stats.Owners = len(owners), len(holders)

	return stats
}

// salesSummary returns paid events after and before the given times (zero - no limit), and their prices.
func salesSummary(events []Event, after, before time.Time) EventAndAveragePrices {
	summary := EventAndAveragePrices{
		Events:       []Event{},
		FloorPrice:   new(big.Float),
		HighestPrice: new(big.Float),
		AveragePrice: new(big.Float),
		Volume:       new(big.Float),
	}
	for _, event := range events {
		if event.Price.Sign() <= 0 ||
			(!after.IsZero() && event.Timestamp.Before(after)) ||
			(!before.IsZero() && event.Timestamp.After(before)) {
			continue
		}

		if len(summary.Events) == 0 || event.Price.Cmp(summary.FloorPrice) < 0 {
			summary.FloorPrice = event.Price
		}
		if event.Price.Cmp(summary.HighestPrice) > 0 {
			summary.HighestPrice = event.Price
		}
		summary.Volume = new(big.Float).Add(summary.Volume, event.Price)
		summary.Events = append(summary.Events, event)
	}

	if len(summary.Events) > 0 {
		summary.AveragePrice = new(big.Float).Quo(summary.Volume, big.NewFloat(float64(len(summary.Events))))
	}

	return summary
}
