package ethsig_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/pkg/ethsig"
)

func TestEthsig(t *testing.T) {
	data1 := map[string]string{
		"addressWallet":   "0xe2B32824733d350845c056CedD73c491FC4C1585",
		"addressContract": "0x0c80417acb4b309725de29b1d950bca974120996",
		"privateKey":      "5aefce0a2d473f59578fa7dee6a122d6509af1e0f79fcbee700dfcfeddabe4cc",
		"signature":       "707fb93c61be8d54c6d1fdf4b83c8642831c480194f7cc93ebdd6fe1ac7474ae63efd077cf6398bf00dc0f7ea96be9f3f9a05dfac1382c4d2f1bb11ec46148491b",
	}

	data2 := map[string]string{
		"addressWallet":   "0xe2B32824733d350845c056CedD73c491FC4C1585",
		"addressContract": "0x02a061be81ee0d7dbbd972bc7edee30b7b102a40",
		"privateKey":      "5aefce0a2d473f59578fa7dee6a122d6509af1e0f79fcbee700dfcfeddabe4cc",
		"signature":       "c2ab46d5981fe2fe951a63240af99a639ca13ea9afa1ee640d418690eae178215571ef9d318cc5b79d2c9e184e036f513aea7edcd6517cca75c4790d9ead45fc1b",
	}
	var tokenID2 int64 = 2

	data3 := map[string]string{
		"addressWallet":   "0xb2cdC7EB2F9d2E629ee97BB91700622A42e688b8",
		"addressContract": "0xde07015be3E663954D514418B4014c3b829D212b",
		"privateKey":      "5aefce0a2d473f59578fa7dee6a122d6509af1e0f79fcbee700dfcfeddabe4cc",
		"signature":       "53f7f1e623364fa4f1bd6e7df67a66edcc06cba01462193d397bbe72fcdba31f04e50d9f65e387be1ea2351110166ad2cb882ccc28be05725e6877645941a3471b",
	}
	var value3 = new(big.Int)
	value3.SetString("5000000000000000000", 10)
	var nonce3 int64 = 0

	data4 := map[string]string{
		"addressWallet":       "0xe2B32824733d350845c056CedD73c491FC4C1585",
		"addressSaleContract": "0x2D1526bad5E59eBc8fc1b2eCBC06dbC8f856Cd2D",
		"addressNFTContract":  "0x674cd2ab624a3eD4f93cbB3d25399cd987B6F564",
		"privateKey":          "5aefce0a2d473f59578fa7dee6a122d6509af1e0f79fcbee700dfcfeddabe4cc",
		"signature":           "ec12bc59894c3964187ec8a0ef15f94ae8a84a1ce943d58d7633ba451ef81aa04689669f8d16e249ddb0afafacec5f231b1a41acecc2884d00ca1982235b9ed71c",
	}
	var tokenID4 int64 = 3
	var value4 = new(big.Int)
	value4.SetString("100000000000000000", 10)

	data5 := map[string]string{
		"contractMethodAddress": "0x095ea7b3",
		"to":                    "0x238f30A8CCF64b502fBA2e0C1b90Ac3c07e3B0a6", // contract.
		"from":                  "0x56f088767D91badc379155290c4205c7b917a36E", // user.
		"domainSeperator":       "0x2739d6640de1503427ab7c5bd20094483387d4f8de3af1aeb1cfbf826f1b5b30",
		"privateKey":            "3b2c3b9eec999beb061fd5b9fc60ae7995e9a81504e4f1c0e852ffc532cd0649",
		"signature":             "ff638d8a00cd4327d547d4405b29d0da9b180296852dbd199330f7c88524229077ea3bdf7366a91407a87cc8c854dd3694900fd1f65bc375306bbb7005d261de1b",
	}
	var value5 = new(big.Int)
	value5.SetString("10000000000", 10)
	var nonce5 int64 = 1

	t.Run("GenerateSignature", func(t *testing.T) {
		privateKeyECDSA, err := crypto.HexToECDSA(data1["privateKey"])
		require.NoError(t, err)

		signature, err := ethsig.SignWallet(
			ethsig.Address(data1["addressWallet"]),
			ethsig.Address(data1["addressContract"]),
			privateKeyECDSA,
		)
		require.NoError(t, err)
		assert.Equal(t, signature, ethsig.Signature(data1["signature"]))
	})

	t.Run("GenerateSignatureWithValue", func(t *testing.T) {
		privateKeyECDSA, err := crypto.HexToECDSA(data2["privateKey"])
		require.NoError(t, err)

		signature, err := ethsig.SignWalletWithValue(
			ethsig.Address(data2["addressWallet"]),
			ethsig.Address(data2["addressContract"]),
			tokenID2,
			privateKeyECDSA,
		)
		require.NoError(t, err)
		assert.Equal(t, signature, ethsig.Signature(data2["signature"]))
	})

	t.Run("GenerateSignatureWithValueAndNonce", func(t *testing.T) {
		privateKeyECDSA, err := crypto.HexToECDSA(data3["privateKey"])
		require.NoError(t, err)

		signature, err := ethsig.SignWalletWithNonce(
			ethsig.Address(data3["addressWallet"]),
			ethsig.Address(data3["addressContract"]),
			value3,
			nonce3,
			privateKeyECDSA,
		)
		require.NoError(t, err)
		assert.Equal(t, signature, ethsig.Signature(data3["signature"]))
	})

	t.Run("GenerateSignatureWithTokenIDAndValue", func(t *testing.T) {
		privateKeyECDSA, err := crypto.HexToECDSA(data4["privateKey"])
		require.NoError(t, err)

		signature, err := ethsig.SignNFTTransaction(
			ethsig.Address(data4["addressWallet"]),
			ethsig.Address(data4["addressSaleContract"]),
			ethsig.Address(data4["addressNFTContract"]),
			tokenID4,
			value4,
			privateKeyECDSA,
		)
		require.NoError(t, err)
		assert.Equal(t, signature, ethsig.Signature(data4["signature"]))
	})

	t.Run("GenerateSignatureForApproveERC20", func(t *testing.T) {
		privateKeyECDSA, err := crypto.HexToECDSA(data5["privateKey"])
		require.NoError(t, err)

		signature, err := ethsig.SignERC20Approval(
			ethsig.Hex(data5["contractMethodAddress"]),
			ethsig.Address(data5["to"]),
			value5,
			nonce5,
			ethsig.Address(data5["from"]),
			ethsig.Hex(data5["domainSeperator"]),
			privateKeyECDSA,
		)
		require.NoError(t, err)
		assert.Equal(t, signature, ethsig.Signature(data5["signature"]))
	})
}
