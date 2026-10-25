package nfts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/collections"
	"generator/internal/nftcontracts"
	"generator/whitelist"
)

// ErrNFTs indicates that there was an error in the service.
var ErrNFTs = errs.Class("nfts service error")

// Service is handling nfts related logic.
//
// architecture: Service
type Service struct {
	config      Config
	collections *collections.Service
	whitelist   *whitelist.Service
}

// NewService is a constructor for nfts service.
func NewService(config Config, collections *collections.Service, whitelist *whitelist.Service) *Service {
	return &Service{
		config:      config,
		collections: collections,
		whitelist:   whitelist,
	}
}

// GetDeployContractTransactionData returns the data of the transaction that deploys the collection contract
// with the given sale terms.
func (service *Service) GetDeployContractTransactionData(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, terms Terms) ([]byte, error) {
	var merkleRootArray [32]byte

	project, err := service.collections.Get(ctx, projectID)
	if err != nil {
		return nil, ErrNFTs.Wrap(err)
	}

	if project.UserID != userID {
		return nil, ErrNFTs.New("this project does not belong to the user")
	}

	contractTerms, err := terms.Resolve(project.Name, project.NFTCount)
	if err != nil {
		return nil, ErrNFTs.Wrap(err)
	}

	merkleRoot, err := service.whitelist.GetMerkleRoot(ctx, projectID)
	if err != nil {
		return nil, ErrNFTs.Wrap(err)
	}

	copy(merkleRootArray[:], merkleRoot)
	// a collection published to IPFS is served from there, otherwise by the platform server.
	fullBaseURI := fmt.Sprintf(service.config.BaseURI, project.ID.String())
	if project.MetadataCID != "" {
		fullBaseURI = fmt.Sprintf("ipfs://%s/", project.MetadataCID)
	}
	data, err := nftcontracts.CreateDeployContractTransactionData(fullBaseURI, merkleRootArray, contractTerms.Name,
		big.NewInt(int64(contractTerms.MaxSupply)), contractTerms.PriceWei, big.NewInt(int64(contractTerms.MaxPerMint)))

	return data, ErrNFTs.Wrap(err)
}
