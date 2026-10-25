package ethsig_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/pkg/ethsig"
)

func TestCrypto(t *testing.T) {
	var value float64 = 0.001

	t.Run("EthToWei", func(t *testing.T) {
		eth, err := ethsig.EthToWei(value)
		require.NoError(t, err)

		ethValue := new(big.Int)
		ethValue.SetString("1000000000000000", 10)
		assert.Equal(t, eth, ethValue)
	})

	t.Run("WeiToEth", func(t *testing.T) {
		weiValue := new(big.Int)
		weiValue.SetString("1000000000000000", 10)

		eth := ethsig.WeiToEth(weiValue)
		assert.Equal(t, eth, value)
	})
}
