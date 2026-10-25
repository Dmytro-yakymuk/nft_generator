package dashboard

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Types of collection events.
const (
	// EventMint is a token created by the collection contract.
	EventMint = "mint"
	// EventSale is a transfer paid with ether, e.g. a purchase on a marketplace.
	EventSale = "sale"
	// EventTransfer is a transfer without payment.
	EventTransfer = "transfer"
	// EventBurn is a token destroyed by the owner.
	EventBurn = "burn"
)

// Period is the default period of the sales chart: all time.
const Period = 0

// PageSize is the number of events on a page of the activity table.
const PageSize = 10

// Currency is the currency of event prices: the native token of the network.
const Currency = "ETH"

// Event describes a transfer of a collection token.
type Event struct {
	From      common.Address `json:"from"`
	To        common.Address `json:"to"`
	TokenID   string         `json:"tokenId"`
	Timestamp time.Time      `json:"timestamp"`
	Type      string         `json:"type"`
	// Price is the ether paid for the token, zero if the transfer was not paid.
	Price       *big.Float  `json:"price"`
	Currency    string      `json:"currency"`
	Image       string      `json:"image"`
	TxHash      common.Hash `json:"txHash"`
	BlockNumber uint64      `json:"blockNumber"`
}

// EventsWithPage describes a page of collection events, newest first.
type EventsWithPage struct {
	Events       []Event `json:"events"`
	PreviousPage string  `json:"previousPage"`
	NextPage     string  `json:"nextPage"`
}

// CollectionStats describes collection statistic entity.
type CollectionStats struct {
	// Items is the number of tokens in circulation: minted and not burned.
	Items int `json:"items"`
	// Owners is the number of addresses that hold at least one token.
	Owners int `json:"owners"`
	// Sales is the number of paid events: paid mints and resales.
	Sales int `json:"sales"`
}

// EventAndAveragePrices describes paid events of a period, oldest first, and their prices.
type EventAndAveragePrices struct {
	Events       []Event    `json:"events"`
	FloorPrice   *big.Float `json:"floorPrice"`
	HighestPrice *big.Float `json:"highestPrice"`
	AveragePrice *big.Float `json:"averagePrice"`
	// Volume is the total ether paid in the period.
	Volume *big.Float `json:"volume"`
}

// Config defines configuration for dashboard.
type Config struct {
	// RPCURL is the node the server reads collection events from;
	// empty means the RPC of the network the web app works with (console.server.chain.rpcUrl).
	RPCURL string `json:"rpcUrl"`
	// LogsBlockRange is the widest block range of one eth_getLogs request the node accepts.
	LogsBlockRange uint64 `json:"logsBlockRange"`
}
