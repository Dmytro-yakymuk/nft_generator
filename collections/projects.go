package collections

import (
	"context"
	"html/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/pkg/filebase"
	"generator/pkg/storage"
)

var (
	// ErrNoCollection indicates that collection does not exist.
	ErrNoCollection = errs.Class("collection does not exist")
	// ErrNoSite indicates that collection mint site does not exists.
	ErrNoSite = errs.Class("collection mint site does not exist")
	// ErrNoGeneratedFiles indicates that generated files does not exists.
	ErrNoGeneratedFiles = errs.Class("generated files does not exists")
	// ErrNoAboutSection indicates that about section does not exist.
	ErrNoAboutSection = errs.Class("about section does not exist")
	// ErrNoHeaderSection indicates that header section does not exist.
	ErrNoHeaderSection = errs.Class("header section does not exist")
	// ErrNoMember indicates that member does not exist.
	ErrNoMember = errs.Class("member does not exist")
	// ErrNoMilestone indicates that milestone does not exist.
	ErrNoMilestone = errs.Class("milestone does not exist")
	// ErrNoSocialNetwork indicates that social network does not exist.
	ErrNoSocialNetwork = errs.Class("social network does not exist")
	// ErrCollectionOwnership indicates that collection does not belong to user.
	ErrCollectionOwnership = errs.Class("this collection does not belong to user")
	// ErrInvalidSlug indicates that the website address does not match SlugPattern.
	ErrInvalidSlug = errs.Class("website address must be 3-63 lowercase letters, digits or hyphens")
	// ErrNotPublished indicates that the collection has not been published to IPFS yet.
	ErrNotPublished = errs.Class("collection is not published to IPFS")
	// ErrPinningDisabled indicates that no pinning service is configured.
	ErrPinningDisabled = errs.Class("no IPFS pinning service is configured")
	// ErrSlugTaken indicates that another collection website already uses the address.
	ErrSlugTaken = errs.Class("website address is already taken")
	// ErrContractNotDeployed indicates that the collection contract has not been deployed yet.
	ErrContractNotDeployed = errs.Class("collection contract is not deployed")
)

// Config defines configurable values of project.
type Config struct {
	ImagesPath  template.URL   `json:"imagesPath"`
	FileStorage storage.Config `json:"fileStorage"`
	Bucket      string         `json:"bucket"`
	// IPFSAPIURL is the RPC API of the IPFS (Kubo) node the collections are published to.
	IPFSAPIURL string `json:"ipfsApiUrl"`
	// Filebase is the account the published collections are pinned in; empty disables pinning.
	Filebase filebase.Config `json:"filebase"`
}

// DB is exposing access to collections db.
//
// architecture: DB
type DB interface {
	// Create creates collection in database.
	Create(ctx context.Context, collection Collection) error
	// Get returns collection by id from database.
	Get(ctx context.Context, id uuid.UUID) (Collection, error)
	// GetByUserID returns all user collections by user id from database.
	GetByUserID(ctx context.Context, id uuid.UUID) ([]Collection, error)
	// List returns all collections from database.
	List(ctx context.Context) ([]Collection, error)
	// Update updates collection in the database.
	Update(ctx context.Context, id uuid.UUID, update Update) error
	// UpdateContractAddress updates contract address and its deployment block for the collection in the database.
	UpdateContractAddress(ctx context.Context, id uuid.UUID, contractAddress common.Address, contractBlock uint64) error
	// UpdateIPFS updates the IPFS directories of the published collection in the database.
	UpdateIPFS(ctx context.Context, id uuid.UUID, imagesCID, metadataCID string) error
	// UpdateIPFSPinned saves the pinning service that keeps the published directories.
	UpdateIPFSPinned(ctx context.Context, id uuid.UUID, service string) error
	// GetBySlug returns collection by its website address from database.
	GetBySlug(ctx context.Context, slug string) (Collection, error)
	// UpdateWebsite updates website address and publication state of the collection in the database.
	UpdateWebsite(ctx context.Context, id uuid.UUID, slug string, isPublished bool) error
	// Delete deletes collection from database.
	Delete(ctx context.Context, id uuid.UUID) error
	// CreateSite creates collection site in the database.
	CreateSite(ctx context.Context, site Site) error
	// DeleteSite deletes collection site from the database.
	DeleteSite(ctx context.Context, id uuid.UUID) error
	// GetHeaderSection returns header section of collection site from the database.
	GetHeaderSection(ctx context.Context, id uuid.UUID) (Section, error)
	// UpdateHeaderSection updates header section of collection site in the database.
	UpdateHeaderSection(ctx context.Context, section Section) error
	// GetAboutSection returns about section of collection site from the database.
	GetAboutSection(ctx context.Context, id uuid.UUID) (Section, error)
	// UpdateAboutSection updates about section of collection site in the database.
	UpdateAboutSection(ctx context.Context, section Section) error
	// AddMember adds member of collection in the database.
	AddMember(ctx context.Context, member Member) error
	// ListMembers returns all collection members from the database.
	ListMembers(ctx context.Context, id uuid.UUID) ([]Member, error)
	// GetMember returns collection member from the database.
	GetMember(ctx context.Context, id uuid.UUID) (Member, error)
	// UpdateMember updates member of collection in the database.
	UpdateMember(ctx context.Context, member Member) error
	// DeleteMember deletes collection member from the database.
	DeleteMember(ctx context.Context, memberID uuid.UUID) error
	// DeleteMembers deletes all collection members from the database.
	DeleteMembers(ctx context.Context, collectionID uuid.UUID) error
	// AddMilestone adds milestone of collection in the database.
	AddMilestone(ctx context.Context, milestone Milestone) error
	// ListMilestones returns all collection milestones from the database.
	ListMilestones(ctx context.Context, id uuid.UUID) ([]Milestone, error)
	// UpdateMilestone updates milestone of collection in the database.
	UpdateMilestone(ctx context.Context, milestone Milestone) error
	// DeleteMilestone deletes collection milestone from the database.
	DeleteMilestone(ctx context.Context, milestoneID uuid.UUID) error
	// DeleteMilestones deletes all collection milestones.
	DeleteMilestones(ctx context.Context, collectionID uuid.UUID) error
	// AddSocialNetwork adds social network for collection site in the database.
	AddSocialNetwork(ctx context.Context, socialNetwork SocialNetwork) error
	// ListSocialNetworks returns social network of collection site from the database.
	ListSocialNetworks(ctx context.Context, id uuid.UUID) ([]SocialNetwork, error)
	// UpdateSocialNetwork updates social network of collection site in the database.
	UpdateSocialNetwork(ctx context.Context, socialNetwork SocialNetwork) error
	// DeleteSocialNetwork updates social network of collection site in the database.
	DeleteSocialNetwork(ctx context.Context, socialNetworkID uuid.UUID) error
	// DeleteSocialNetworks deletes all collection social networks from the database.
	DeleteSocialNetworks(ctx context.Context, collectionID uuid.UUID) error
}

// Collection describes collection entity.
type Collection struct {
	ID              uuid.UUID      `json:"id"`
	Name            string         `json:"name"`
	NFTCount        int            `json:"nftCount"`
	UserID          uuid.UUID      `json:"userId"`
	ContractAddress common.Address `json:"contractAddress"`
	// ContractBlock is the number of the block the contract was deployed in:
	// the dashboard reads contract events starting from it.
	ContractBlock      uint64 `json:"contractBlock"`
	AreImagesGenerated bool   `json:"areImagesGenerated"`
	IsWebsitePublished bool   `json:"isWebsitePublished"`
	// WebsiteSlug is the address of the published mint site: /sites/{WebsiteSlug}.
	WebsiteSlug string `json:"websiteSlug"`
	// ImagesCID and MetadataCID are the IPFS directories of the published collection:
	// <tokenId>.png images and <tokenId> metadata files; the contract base URI is ipfs://MetadataCID/.
	ImagesCID   string `json:"imagesCid"`
	MetadataCID string `json:"metadataCid"`
	// IPFSPinnedIn names the pinning service that keeps both directories, empty while only the
	// platform's node holds them.
	IPFSPinnedIn string `json:"ipfsPinnedIn"`
}

// Update defines required fields for project updating.
type Update struct {
	Name     string `json:"name"`
	NFTCount int    `json:"nftCount"`
}

// Section defines required fields for header and about section of project site.
type Section struct {
	ProjectID   uuid.UUID `json:"projectId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

// Member defines required fields for project team section of project site.
type Member struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
}

// Milestone defines required fields for project roadmap of project site.
type Milestone struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"projectId"`
	Title       string    `json:"title"`
	Date        string    `json:"date"`
	Description string    `json:"description"`
}

// SocialNetwork defines required fields for project social networks of project site.
type SocialNetwork struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	Title     string    `json:"title"`
	Link      string    `json:"link"`
}

// Site describes representation of project site.
type Site struct {
	ProjectID      uuid.UUID       `json:"projectId"`
	Name           string          `json:"name"`
	About          Section         `json:"aboutSection"`
	Header         Section         `json:"headerSection"`
	ContentNames   []string        `json:"contentNames"`
	Team           []Member        `json:"projectMemberSections"`
	Roadmap        []Milestone     `json:"roadmap"`
	SocialNetworks []SocialNetwork `json:"socialNetworks"`
	ImagesPath     template.URL    `json:"imagesPath"`
}
