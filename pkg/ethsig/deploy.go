package ethsig

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeebo/errs"
)

// ErrDeployContract indicates that an error occurred while deploying a contract.
var ErrDeployContract = errs.Class("deploy contract error")

// ContractConfig describes values required to deploy contact.
type ContractConfig struct {
	AddressNodeServer string `json:"addressNodeServer"`
	PrivateKey        string `json:"privateKey"`
}

// DeployNFTContract deploys nft contract for each collection nfts.
func DeployNFTContract(ctx context.Context, config ContractConfig, deployContract func(auth *bind.TransactOpts, client *ethclient.Client) (common.Address, error)) (common.Address, error) {
	client, err := ethclient.Dial(config.AddressNodeServer)
	if err != nil {
		return common.Address{}, ErrDeployContract.Wrap(err)
	}

	privateKeyECDSA, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return common.Address{}, ErrDeployContract.Wrap(err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return common.Address{}, ErrDeployContract.Wrap(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)
	if err != nil {
		return common.Address{}, ErrDeployContract.Wrap(err)
	}
	auth.Value = big.NewInt(0) // in wei.

	address, err := deployContract(auth, client)
	return address, ErrDeployContract.Wrap(err)
}
