package nftcontracts_test

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"generator/internal/nftcontracts"
)

var (
	chainID     = flag.String("chain-id", getenv("CHAIN_ID"), "empty chain id")
	nodeAddress = flag.String("node-address", getenv("NODE_ADDRESS"), "empty node address")
	privateKey  = flag.String("private-key", getenv("PRIVATE_KEY"), "empty private key")
)

func getenv(priority ...string) string {
	for _, p := range priority {
		v := os.Getenv(p)
		if v != "" {
			return v
		}
	}
	return ""
}

func TestNFTContracts(t *testing.T) {
	baseURI := "https://example.com/nfts/"
	merkleRoot := [32]byte{1, 2, 3}
	name := "Test Contract"
	maxSupply, price, maxPerMint := big.NewInt(100), big.NewInt(1e15), big.NewInt(5)

	t.Run("CreateDeployContractTransactionData", func(t *testing.T) {
		data, err := nftcontracts.CreateDeployContractTransactionData(baseURI, merkleRoot, name, maxSupply, price, maxPerMint)
		require.NoError(t, err)

		// the bytecode is followed by the constructor arguments in their order.
		contractABI, err := abi.JSON(strings.NewReader(nftcontracts.CollectionNFTMetaData.ABI))
		require.NoError(t, err)
		bin := common.FromHex(nftcontracts.CollectionNFTMetaData.Bin)
		require.Equal(t, bin, data[:len(bin)])
		args, err := contractABI.Constructor.Inputs.Unpack(data[len(bin):])
		require.NoError(t, err)
		require.Equal(t, []interface{}{baseURI, merkleRoot, name, maxSupply, price, maxPerMint}, args)
	})

	t.Run("DeployContract", func(t *testing.T) {
		checkFlag(t)

		privateKeyECDSA, err := crypto.HexToECDSA(*privateKey)
		require.NoError(t, err)

		publicKey := privateKeyECDSA.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}

		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		client, err := ethclient.Dial(*nodeAddress)
		require.NoError(t, err)

		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		require.NoError(t, err)
		value := big.NewInt(0)
		gasLimit := uint64(4300000)
		require.NoError(t, err)

		data, err := nftcontracts.CreateDeployContractTransactionData(baseURI, merkleRoot, name, maxSupply, price, maxPerMint)
		require.NoError(t, err)

		chainIDInt, err := strconv.ParseInt(*chainID, 10, 64)
		require.NoError(t, err)

		tx := types.NewTx(&types.LegacyTx{
			Nonce: nonce,
			Gas:   gasLimit,
			Value: value,
			Data:  data,
		})

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainIDInt)), privateKeyECDSA)
		require.NoError(t, err)

		err = client.SendTransaction(context.Background(), signedTx)
		require.NoError(t, err)
	})
}

func checkFlag(t *testing.T) {
	if *privateKey == "" {
		t.Skip("skip testing because private key is empty")
	}

	if *chainID == "" {
		*chainID = "5777"
	}

	if *nodeAddress == "" {
		*nodeAddress = "http://127.0.0.1:7545"
	}
}
