package nftcontracts

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/zeebo/errs"

	"generator/pkg/ethsig"
)

// ErrorCreateData indicates that an error occurs when creating a date.
var ErrorCreateData = errs.Class("create data error")

// CreateDeployContractTransactionData creates and returns data for deploying contract: the bytecode followed
// by the constructor arguments (base URI, whitelist Merkle root, name and the sale terms, the price in wei).
func CreateDeployContractTransactionData(baseURI string, merkleRoot [32]byte, name string, maxSupply, price, maxPerMint *big.Int) ([]byte, error) {
	var data []byte

	abiReader := strings.NewReader(CollectionNFTMetaData.ABI)
	abi, err := abi.JSON(abiReader)
	if err != nil {
		return nil, ErrorCreateData.Wrap(err)
	}

	constructor, err := abi.Pack("", baseURI, merkleRoot, name, maxSupply, price, maxPerMint)
	if err != nil {
		return nil, ErrorCreateData.Wrap(err)
	}

	bin, err := hex.DecodeString(CollectionNFTMetaData.Bin[ethsig.LengthHexPrefix:])
	if err != nil {
		return nil, ErrorCreateData.Wrap(err)
	}

	data = append(data, append(bin, constructor...)...)
	return data, nil
}
