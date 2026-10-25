package nfts

import (
	"math/big"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/zeebo/errs"
)

// Config defines values required to deploy contract.
type Config struct {
	BaseURI string `json:"baseUri"`
}

// Default sale terms: the contract page and the command line start from them, and the mint site falls back
// to them for contracts deployed before the terms became parameters of the constructor.
const (
	// DefaultMintPrice is the price of one token, ETH.
	DefaultMintPrice = "0.001"
	// DefaultMaxPerMint is the maximum number of tokens minted in one transaction.
	DefaultMaxPerMint = 5
)

// MaxNameLength limits the contract name, which is stored on chain.
const MaxNameLength = 64

// ErrInvalidTerms indicates sale terms the contract cannot be deployed with.
var ErrInvalidTerms = errs.Class("invalid sale terms")

// etherAmount is a non-negative amount of ETH that converts into wei exactly.
var etherAmount = regexp.MustCompile(`^[0-9]{1,9}(\.[0-9]{1,18})?$`)

// Terms are the sale terms the owner chooses for the collection contract; empty fields take the defaults.
type Terms struct {
	// Name of the contract (ERC-721 name); the collection name by default.
	Name string
	// Price of one token, ETH, as a decimal string.
	Price string
	// MaxPerMint is the limit of tokens per transaction.
	MaxPerMint int
}

// ContractTerms are the checked constructor arguments of the collection contract.
type ContractTerms struct {
	Name       string
	MaxSupply  int
	PriceWei   *big.Int
	MaxPerMint int
}

// Resolve fills the defaults and checks the terms against the collection: the contract can mint exactly
// as many tokens as the collection has metadata for.
func (terms Terms) Resolve(collectionName string, nftCount int) (ContractTerms, error) {
	if nftCount <= 0 {
		return ContractTerms{}, ErrInvalidTerms.New("the collection has no tokens")
	}

	name := strings.TrimSpace(terms.Name)
	if name == "" {
		name = collectionName
	}
	if utf8.RuneCountInString(name) > MaxNameLength {
		return ContractTerms{}, ErrInvalidTerms.New("the contract name is longer than %d characters", MaxNameLength)
	}

	price := strings.TrimSpace(terms.Price)
	if price == "" {
		price = DefaultMintPrice
	}
	priceWei, err := ParseEther(price)
	if err != nil {
		return ContractTerms{}, err
	}

	maxPerMint := terms.MaxPerMint
	if maxPerMint == 0 {
		maxPerMint = DefaultMaxPerMint
		if maxPerMint > nftCount {
			maxPerMint = nftCount
		}
	}
	if maxPerMint < 1 || maxPerMint > nftCount {
		return ContractTerms{}, ErrInvalidTerms.New("the limit per transaction must be from 1 to %d", nftCount)
	}

	return ContractTerms{Name: name, MaxSupply: nftCount, PriceWei: priceWei, MaxPerMint: maxPerMint}, nil
}

// ParseEther converts a decimal amount of ETH into wei without rounding.
func ParseEther(amount string) (*big.Int, error) {
	if !etherAmount.MatchString(amount) {
		return nil, ErrInvalidTerms.New("the price must be a number of ETH with at most 18 decimals, got %q", amount)
	}

	parts := strings.SplitN(amount, ".", 2)
	fraction := ""
	if len(parts) == 2 {
		fraction = parts[1]
	}
	wei, _ := new(big.Int).SetString(parts[0]+fraction+strings.Repeat("0", 18-len(fraction)), 10)

	return wei, nil
}
