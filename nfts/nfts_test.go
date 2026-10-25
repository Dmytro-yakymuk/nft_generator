package nfts_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"generator/nfts"
)

func TestParseEther(t *testing.T) {
	for amount, wei := range map[string]string{
		"0":                    "0",
		"0.001":                "1000000000000000",
		"1":                    "1000000000000000000",
		"2.5":                  "2500000000000000000",
		"0.000000000000000001": "1",
		"123456789.1":          "123456789100000000000000000",
	} {
		got, err := nfts.ParseEther(amount)
		require.NoError(t, err, amount)
		require.Equal(t, wei, got.String(), amount)
	}

	for _, amount := range []string{"", "-1", "1.", ".5", "1e3", "0.0000000000000000001", "1,5", "1234567890", " 1"} {
		_, err := nfts.ParseEther(amount)
		require.Error(t, err, amount)
		require.True(t, nfts.ErrInvalidTerms.Has(err), amount)
	}
}

func TestTermsResolve(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		terms, err := nfts.Terms{}.Resolve("Cosmic Critters", 100)
		require.NoError(t, err)
		require.Equal(t, nfts.ContractTerms{Name: "Cosmic Critters", MaxSupply: 100, PriceWei: big.NewInt(1e15), MaxPerMint: nfts.DefaultMaxPerMint}, terms)
	})

	t.Run("the default limit fits a small collection", func(t *testing.T) {
		terms, err := nfts.Terms{}.Resolve("Small", 3)
		require.NoError(t, err)
		require.Equal(t, 3, terms.MaxPerMint)
	})

	t.Run("chosen terms", func(t *testing.T) {
		terms, err := nfts.Terms{Name: "  Critters  ", Price: "0.02", MaxPerMint: 10}.Resolve("Cosmic Critters", 100)
		require.NoError(t, err)
		require.Equal(t, nfts.ContractTerms{Name: "Critters", MaxSupply: 100, PriceWei: big.NewInt(2e16), MaxPerMint: 10}, terms)
	})

	t.Run("free mint", func(t *testing.T) {
		terms, err := nfts.Terms{Price: "0"}.Resolve("Free", 10)
		require.NoError(t, err)
		require.Zero(t, terms.PriceWei.Sign())
	})

	for name, tc := range map[string]struct {
		terms    nfts.Terms
		nftCount int
	}{
		"no tokens":            {nfts.Terms{}, 0},
		"limit above the size": {nfts.Terms{MaxPerMint: 11}, 10},
		"negative limit":       {nfts.Terms{MaxPerMint: -1}, 10},
		"bad price":            {nfts.Terms{Price: "free"}, 10},
		"long name":            {nfts.Terms{Name: strings.Repeat("я", nfts.MaxNameLength+1)}, 10},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			_, err := tc.terms.Resolve("Collection", tc.nftCount)
			require.True(t, nfts.ErrInvalidTerms.Has(err), err)
		})
	}
}
