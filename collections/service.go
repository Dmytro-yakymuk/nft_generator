package collections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"generator/collections/images"
	"generator/collections/layers"
	"generator/collections/metadata"
	"generator/internal/filestore"
	"generator/internal/filestore/localfilestore"
	"generator/internal/nft"
	"generator/pkg/filebase"
	"generator/pkg/ipfs"
)

// ErrCollections indicates that there was an error in the service.
var ErrCollections = errs.Class("collections service error")

// ErrNoToken indicates that the collection has no token with the given id.
var ErrNoToken = errs.Class("token does not exist")

// Service is handling projects related logic.
//
// architecture: Service
type Service struct {
	config   Config
	projects DB
	store    filestore.FileStore
	layers   *layers.Service
	images   *images.Service
	metadata *metadata.Service
	ipfs     *ipfs.Client
	// pinning keeps the published directories in a pinning service; nil if none is configured.
	pinning *filebase.Client
}

// NewService is a constructor for projects service.
func NewService(cfg Config, projects DB, layers *layers.Service, fileStore filestore.FileStore, images *images.Service, metadata *metadata.Service) *Service {
	service := &Service{
		config:   cfg,
		projects: projects,
		store:    fileStore,
		layers:   layers,
		images:   images,
		metadata: metadata,
		ipfs:     ipfs.NewClient(cfg.IPFSAPIURL),
	}
	if cfg.Filebase.Enabled() {
		service.pinning = filebase.NewClient(cfg.Filebase)
	}

	return service
}

// PinningService returns the name of the configured pinning service, or an empty string.
func (service *Service) PinningService() string {
	if service.pinning == nil {
		return ""
	}

	return filebase.Service
}

// Create creates new collection.
func (service *Service) Create(ctx context.Context, project Collection) error {
	return ErrCollections.Wrap(service.projects.Create(ctx, project))
}

// List returns user collections.
func (service *Service) List(ctx context.Context, userID uuid.UUID) ([]Collection, error) {
	userProjects, err := service.projects.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrCollections.Wrap(err)
	}
	for i := range userProjects {
		service.setImagesGenerated(ctx, &userProjects[i])
	}
	return userProjects, nil
}

// Get returns collection by id.
func (service *Service) Get(ctx context.Context, id uuid.UUID) (Collection, error) {
	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return project, ErrCollections.Wrap(err)
	}
	service.setImagesGenerated(ctx, &project)
	return project, nil
}

// setImagesGenerated marks the collection as generated when all its images exist. The flag is derived
// from the generated files rather than stored: a running or cancelled generation leaves fewer images,
// and a new generation deletes the old ones first.
func (service *Service) setImagesGenerated(ctx context.Context, project *Collection) {
	if service.images == nil {
		return
	}
	count, err := service.images.GeneratedImagesCount(ctx, project.ID)
	project.AreImagesGenerated = err == nil && project.NFTCount > 0 && count >= project.NFTCount
}

// SlugPattern describes the mint site address: 3-63 lowercase letters, digits or hyphens
// that start and end with a letter or digit, so the address can also serve as a subdomain.
var SlugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$`)

// PublishSite publishes the collection mint site at /sites/{slug}. The site needs a deployed
// contract to mint from and a filled header section to show.
func (service *Service) PublishSite(ctx context.Context, id, userID uuid.UUID, slug string) (Collection, error) {
	if !SlugPattern.MatchString(slug) {
		return Collection{}, ErrInvalidSlug.New("%q", slug)
	}

	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return Collection{}, ErrCollections.Wrap(err)
	}
	if project.UserID != userID {
		return Collection{}, ErrCollectionOwnership.New("")
	}
	if project.ContractAddress == (common.Address{}) {
		return Collection{}, ErrContractNotDeployed.New("deploy the contract before publishing the website")
	}
	if _, err = service.projects.GetHeaderSection(ctx, id); err != nil {
		return Collection{}, ErrCollections.Wrap(err)
	}

	if err = service.projects.UpdateWebsite(ctx, id, slug, true); err != nil {
		return Collection{}, ErrCollections.Wrap(err)
	}
	project.WebsiteSlug, project.IsWebsitePublished = slug, true

	return project, nil
}

// UnpublishSite hides the collection mint site; the address stays reserved for the collection.
func (service *Service) UnpublishSite(ctx context.Context, id, userID uuid.UUID) error {
	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}
	if project.UserID != userID {
		return ErrCollectionOwnership.New("")
	}

	return ErrCollections.Wrap(service.projects.UpdateWebsite(ctx, id, project.WebsiteSlug, false))
}

// GetPublishedSite returns the published mint site and its collection by the website address.
func (service *Service) GetPublishedSite(ctx context.Context, slug string) (Site, Collection, error) {
	project, err := service.projects.GetBySlug(ctx, slug)
	if err != nil {
		return Site{}, Collection{}, ErrCollections.Wrap(err)
	}
	if !project.IsWebsitePublished {
		return Site{}, Collection{}, ErrNoCollection.New("website %q is not published", slug)
	}

	site, err := service.GetSite(ctx, project.ID)

	return site, project, err
}

// GetSite returns data representation of collection site.
func (service *Service) GetSite(ctx context.Context, id uuid.UUID) (Site, error) {
	var site Site

	project, err := service.Get(ctx, id)
	if err != nil {
		return site, err
	}

	site.ProjectID = id
	site.Name = project.Name
	site.ImagesPath = service.config.ImagesPath

	site.Header, err = service.GetHeaderSection(ctx, id)
	if err != nil {
		return site, err
	}

	site.About, err = service.GetAboutSection(ctx, id)
	if err != nil {
		return site, err
	}

	site.ContentNames, err = service.ListContentFiles(ctx, id)
	if err != nil && !localfilestore.ErrNoFile.Has(err) {
		return site, err
	}

	site.Roadmap, err = service.ListMilestones(ctx, id)
	if err != nil {
		return site, err
	}

	site.Team, err = service.ListMembers(ctx, id)
	if err != nil {
		return site, err
	}

	site.SocialNetworks, err = service.ListSocialNetworks(ctx, id)

	return site, err
}

// DeleteSite deletes collection site and all data related to it.
func (service *Service) DeleteSite(ctx context.Context, id uuid.UUID) error {
	err := service.projects.DeleteSite(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	err = service.DeleteLogo(ctx, id)
	e, ok := err.(*os.PathError)
	if err != nil && ok && errs.Is(e.Err, syscall.ENOENT) {
		return err
	}

	err = service.DeleteHeaderSectionImage(ctx, id)
	e, ok = err.(*os.PathError)
	if err != nil && ok && errs.Is(e.Err, syscall.ENOENT) {
		return err
	}

	err = service.DeleteAboutSectionImage(ctx, id)
	e, ok = err.(*os.PathError)
	if err != nil && ok && errs.Is(e.Err, syscall.ENOENT) {
		return err
	}

	members, err := service.ListMembers(ctx, id)
	if err != nil {
		return err
	}

	for _, member := range members {
		err = service.DeleteMemberImage(ctx, id, member.ID)
		e, ok = err.(*os.PathError)
		if err != nil && ok && errs.Is(e.Err, syscall.ENOENT) {
			return err
		}
	}

	err = service.projects.DeleteMembers(ctx, id)
	if err != nil && !ErrNoMember.Has(err) {
		return ErrCollections.Wrap(err)
	}

	contentImages, err := service.ListContentFiles(ctx, id)
	e, ok = err.(*os.PathError)
	if err != nil && ok && errs.Is(e.Err, syscall.ENOENT) {
		return err
	}

	for _, name := range contentImages {
		err = service.DeleteContentImage(ctx, id, name)
		if err != nil {
			return ErrCollections.Wrap(err)
		}
	}

	err = service.projects.DeleteMilestones(ctx, id)
	if err != nil && !ErrNoMilestone.Has(err) {
		return ErrCollections.Wrap(err)
	}

	err = service.projects.DeleteSocialNetworks(ctx, id)
	if err != nil && !ErrNoSocialNetwork.Has(err) {
		return ErrCollections.Wrap(err)
	}

	return nil
}

// Generate generates collection nfts from layers.
func (service *Service) Generate(ctx context.Context, id uuid.UUID) error {
	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	if err = service.images.DeleteGeneratedImages(ctx, project.ID); err != nil {
		return ErrCollections.Wrap(err)
	}

	if err = service.metadata.Delete(ctx, project.ID); err != nil {
		return ErrCollections.Wrap(err)
	}

	return ErrCollections.Wrap(service.layers.Generate(ctx, id, project.NFTCount))
}

// GenerateMetadata generates metadata of nfts for collections.
func (service *Service) GenerateMetadata(ctx context.Context, id uuid.UUID) error {
	_, err := service.projects.Get(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	images, err := service.images.List(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	for _, image := range images {
		number, ok := new(big.Int).SetString(strings.Split(image, ".")[0], 10)
		if !ok {
			return ErrCollections.New("image %q is not named by a number", image)
		}

		// the combination is kept in the dna of the metadata; collections generated before had
		// the images named by the combination number itself.
		combination := number
		if metadata, err := service.metadata.Get(ctx, id, image); err == nil && metadata.DNA != "" {
			if dna, ok := new(big.Int).SetString(metadata.DNA, 10); ok {
				combination = dna
			}
		}

		attributes, err := service.layers.GetImageComponents(ctx, id, combination)
		if err != nil {
			return ErrCollections.Wrap(err)
		}

		if err = service.metadata.GenerateMetadata(ctx, id, number, combination.String(), attributes); err != nil {
			return ErrCollections.Wrap(err)
		}
	}

	return nil
}

// TokenMetadata returns metadata of the token with the given id. Token ids are
// assigned to generated images in the same order as in UploadFiles, so the
// served metadata matches the files uploaded to the S3 storage.
func (service *Service) TokenMetadata(ctx context.Context, id uuid.UUID, tokenID int) (nft.Metadata, error) {
	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return nft.Metadata{}, ErrCollections.Wrap(err)
	}

	images, err := service.images.List(ctx, id)
	if err != nil && !localfilestore.ErrNoFile.Has(err) {
		return nft.Metadata{}, ErrCollections.Wrap(err)
	}

	if tokenID < 0 || tokenID >= len(images) {
		return nft.Metadata{}, ErrNoToken.New("token %d of collection %s", tokenID, id)
	}

	metadata, err := service.metadata.Get(ctx, id, images[tokenID])
	if err != nil {
		return nft.Metadata{}, ErrCollections.Wrap(err)
	}

	metadata.Name = fmt.Sprintf("%s #%d", project.Name, tokenID)

	return metadata, nil
}

// PublishToIPFS adds the generated collection to IPFS as two directories: images <tokenId>.png and
// metadata files <tokenId> whose image is ipfs://<images CID>/<tokenId>.png. Token ids follow the order
// of TokenMetadata, so with the base URI ipfs://<metadata CID>/ the contract's tokenURI(n) returns the
// same metadata as the platform, but no longer depends on the platform server.
func (service *Service) PublishToIPFS(ctx context.Context, id uuid.UUID) (Collection, error) {
	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return project, ErrCollections.Wrap(err)
	}

	names, err := service.images.List(ctx, id)
	if err != nil && !localfilestore.ErrNoFile.Has(err) {
		return project, ErrCollections.Wrap(err)
	}
	if len(names) == 0 {
		return project, ErrNoGeneratedFiles.New("collection %s", id)
	}

	imageFiles := make([]ipfs.File, len(names))
	for tokenID, name := range names {
		name := name
		imageFiles[tokenID] = ipfs.File{
			Name:    fmt.Sprintf("%d.png", tokenID),
			Content: func() ([]byte, error) { return service.images.Get(ctx, id, name) },
		}
	}
	imagesCID, err := service.ipfs.AddDirectory(ctx, imageFiles)
	if err != nil {
		return project, ErrCollections.Wrap(err)
	}

	metadataFiles := make([]ipfs.File, len(names))
	for tokenID, name := range names {
		tokenID, name := tokenID, name
		metadataFiles[tokenID] = ipfs.File{
			Name: strconv.Itoa(tokenID),
			Content: func() ([]byte, error) {
				metadata, err := service.metadata.Get(ctx, id, name)
				if err != nil {
					return nil, err
				}
				metadata.Name = fmt.Sprintf("%s #%d", project.Name, tokenID)
				metadata.Image = fmt.Sprintf("ipfs://%s/%d.png", imagesCID, tokenID)
				return json.Marshal(metadata)
			},
		}
	}
	metadataCID, err := service.ipfs.AddDirectory(ctx, metadataFiles)
	if err != nil {
		return project, ErrCollections.Wrap(err)
	}

	if err = service.projects.UpdateIPFS(ctx, id, imagesCID, metadataCID); err != nil {
		return project, ErrCollections.Wrap(err)
	}
	project.ImagesCID, project.MetadataCID = imagesCID, metadataCID

	return project, nil
}

// PinIPFS pins both published directories of the collection in the pinning service: the platform's node
// exports each as a CAR archive, which the service imports under the same CID, so the contract's
// ipfs:// base URI keeps working without the platform. Pinning the same directories again is harmless.
func (service *Service) PinIPFS(ctx context.Context, id uuid.UUID) (Collection, error) {
	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return project, ErrCollections.Wrap(err)
	}
	if service.pinning == nil {
		return project, ErrPinningDisabled.New("")
	}
	if project.ImagesCID == "" || project.MetadataCID == "" {
		return project, ErrNotPublished.New("collection %s", id)
	}

	for _, cid := range []string{project.ImagesCID, project.MetadataCID} {
		if err = service.pinCAR(ctx, cid); err != nil {
			return project, ErrCollections.Wrap(err)
		}
	}

	if err = service.projects.UpdateIPFSPinned(ctx, id, filebase.Service); err != nil {
		return project, ErrCollections.Wrap(err)
	}
	project.IPFSPinnedIn = filebase.Service

	return project, nil
}

// pinCAR exports the directory to a temporary CAR file and uploads it to the pinning service.
func (service *Service) pinCAR(ctx context.Context, cid string) (err error) {
	car, err := os.CreateTemp("", "ipfs-*.car")
	if err != nil {
		return err
	}
	defer func() {
		err = errs.Combine(err, car.Close(), os.Remove(car.Name()))
	}()

	if err = service.ipfs.ExportCAR(ctx, cid, car); err != nil {
		return err
	}
	if _, err = car.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return service.pinning.PinCAR(ctx, cid, car)
}

// UploadFiles uploads files for collections.
func (service *Service) UploadFiles(ctx context.Context, id uuid.UUID) error {
	_, err := service.projects.Get(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	images, err := service.images.List(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	for k, image := range images {
		if err = service.images.Upload(ctx, id, image); err != nil {
			return ErrCollections.Wrap(err)
		}

		if err = service.metadata.Upload(ctx, id, image, strconv.Itoa(k)); err != nil {
			return ErrCollections.Wrap(err)
		}
	}

	return nil
}

// Update updates certain collection.
func (service *Service) Update(ctx context.Context, id uuid.UUID, update Update) error {
	return ErrCollections.Wrap(service.projects.Update(ctx, id, update))
}

// UpdateContractAddress saves the address of the deployed collection contract and the block it was deployed in.
func (service *Service) UpdateContractAddress(ctx context.Context, id, userID uuid.UUID, contractAddress common.Address, contractBlock uint64) error {
	project, err := service.projects.Get(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}
	if project.UserID != userID {
		return ErrCollectionOwnership.New("")
	}

	return ErrCollections.Wrap(service.projects.UpdateContractAddress(ctx, id, contractAddress, contractBlock))
}

// Delete deletes certain collection.
func (service *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := service.projects.Delete(ctx, id)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	return ErrCollections.Wrap(service.images.DeleteProjectImages(ctx, id))
}

// UpdateHeaderSection updates header section.
func (service *Service) UpdateHeaderSection(ctx context.Context, section Section) error {
	_, err := service.GetHeaderSection(ctx, section.ProjectID)
	if err != nil {
		if ErrNoHeaderSection.Has(err) {
			var site Site
			site.Header = section

			return ErrCollections.Wrap(service.projects.CreateSite(ctx, site))
		}

		return ErrCollections.Wrap(err)
	}

	return ErrCollections.Wrap(service.projects.UpdateHeaderSection(ctx, section))
}

// UpdateHeaderSectionImage updates image for collection header section.
func (service *Service) UpdateHeaderSectionImage(ctx context.Context, projectID uuid.UUID, reader io.Reader) error {
	_, err := service.projects.GetHeaderSection(ctx, projectID)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	relatedPath := filepath.Join(projectID.String(), "header.png")

	return ErrCollections.Wrap(service.store.Store(ctx, relatedPath, reader))
}

// GetHeaderSection returns header section of certain collection.
func (service *Service) GetHeaderSection(ctx context.Context, id uuid.UUID) (Section, error) {
	section, err := service.projects.GetHeaderSection(ctx, id)

	return section, ErrCollections.Wrap(err)
}

// DeleteHeaderSectionImage deletes collection header section image.
func (service *Service) DeleteHeaderSectionImage(ctx context.Context, id uuid.UUID) error {
	path := filepath.Join(id.String(), "header.png")

	return ErrCollections.Wrap(service.store.Delete(ctx, path))
}

// UpdateAboutSection updates about section.
func (service *Service) UpdateAboutSection(ctx context.Context, section Section) error {
	_, err := service.GetAboutSection(ctx, section.ProjectID)
	if err != nil {
		if ErrNoAboutSection.Has(err) {
			var site Site
			site.About = section

			return ErrCollections.Wrap(service.projects.CreateSite(ctx, site))
		}

		return ErrCollections.Wrap(err)
	}

	return ErrCollections.Wrap(service.projects.UpdateAboutSection(ctx, section))
}

// UpdateLogo updates logo of collection.
func (service *Service) UpdateLogo(ctx context.Context, projectID uuid.UUID, reader io.Reader) error {
	relatedPath := filepath.Join(projectID.String(), "logo.png")

	return ErrCollections.Wrap(service.store.Store(ctx, relatedPath, reader))
}

// DeleteLogo deletes logo of collection.
func (service *Service) DeleteLogo(ctx context.Context, projectID uuid.UUID) error {
	relatedPath := filepath.Join(projectID.String(), "logo.png")

	return ErrCollections.Wrap(service.store.Delete(ctx, relatedPath))
}

// UpdateAboutSectionImage updates image for collection about section.
func (service *Service) UpdateAboutSectionImage(ctx context.Context, projectID uuid.UUID, reader io.Reader) error {
	_, err := service.projects.GetAboutSection(ctx, projectID)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	relatedPath := filepath.Join(projectID.String(), "about.png")

	return ErrCollections.Wrap(service.store.Store(ctx, relatedPath, reader))
}

// GetAboutSection returns about section of certain collection.
func (service *Service) GetAboutSection(ctx context.Context, id uuid.UUID) (Section, error) {
	section, err := service.projects.GetAboutSection(ctx, id)

	return section, ErrCollections.Wrap(err)
}

// DeleteAboutSectionImage deletes collection about section image.
func (service *Service) DeleteAboutSectionImage(ctx context.Context, id uuid.UUID) error {
	path := filepath.Join(id.String(), "about.png")

	return ErrCollections.Wrap(service.store.Delete(ctx, path))
}

// AddMember adds new member.
func (service *Service) AddMember(ctx context.Context, projectID uuid.UUID, name, role string) error {
	member := Member{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      name,
		Role:      role,
	}

	return ErrCollections.Wrap(service.projects.AddMember(ctx, member))
}

// UpdateMemberImage adds image for member.
func (service *Service) UpdateMemberImage(ctx context.Context, memberID uuid.UUID, reader io.Reader) error {
	member, err := service.projects.GetMember(ctx, memberID)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	relatedPath := filepath.Join(member.ProjectID.String(), "members", member.ID.String()+".png")

	return ErrCollections.Wrap(service.store.Store(ctx, relatedPath, reader))
}

// ListMembers returns collection members.
func (service *Service) ListMembers(ctx context.Context, id uuid.UUID) ([]Member, error) {
	members, err := service.projects.ListMembers(ctx, id)

	return members, ErrCollections.Wrap(err)
}

// UpdateMember updates member.
func (service *Service) UpdateMember(ctx context.Context, member Member) error {
	return ErrCollections.Wrap(service.projects.UpdateMember(ctx, member))
}

// DeleteMember deletes certain member and deletes member image.
func (service *Service) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	member, err := service.projects.GetMember(ctx, memberID)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	return ErrCollections.Wrap(service.projects.DeleteMember(ctx, member.ID))
}

// DeleteMemberImage deletes member image.
func (service *Service) DeleteMemberImage(ctx context.Context, projectID, memberID uuid.UUID) error {
	_, err := service.projects.GetMember(ctx, memberID)
	if err != nil {
		return ErrCollections.Wrap(err)
	}

	path := filepath.Join(projectID.String(), "members", memberID.String()+".png")

	// a member without a photo is what the request asks for, so a missing file is not an error.
	if err = service.store.Delete(ctx, path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return ErrCollections.Wrap(err)
	}
	return nil
}

// UpdateContentImage adds content image.
func (service *Service) UpdateContentImage(ctx context.Context, projectID uuid.UUID, fileName string, reader io.Reader) error {
	relatedPath := filepath.Join(projectID.String(), "content", fileName)

	return ErrCollections.Wrap(service.store.Store(ctx, relatedPath, reader))
}

// ListContentFiles returns names of content images.
func (service *Service) ListContentFiles(ctx context.Context, projectID uuid.UUID) ([]string, error) {
	path := filepath.Join(projectID.String(), "content")
	names, err := service.store.List(ctx, path)

	return names, ErrCollections.Wrap(err)
}

// DeleteContentImage deletes content image.
func (service *Service) DeleteContentImage(ctx context.Context, projectID uuid.UUID, fileName string) error {
	path := filepath.Join(projectID.String(), "content", fileName)

	return ErrCollections.Wrap(service.store.Delete(ctx, path))
}

// AddMilestone adds new milestone.
func (service *Service) AddMilestone(ctx context.Context, projectID uuid.UUID, date, title, description string) error {
	milestone := Milestone{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Date:        date,
		Title:       title,
		Description: description,
	}

	return ErrCollections.Wrap(service.projects.AddMilestone(ctx, milestone))
}

// ListMilestones returns collection milestones.
func (service *Service) ListMilestones(ctx context.Context, id uuid.UUID) ([]Milestone, error) {
	milestones, err := service.projects.ListMilestones(ctx, id)

	return milestones, ErrCollections.Wrap(err)
}

// UpdateMilestone updates milestone.
func (service *Service) UpdateMilestone(ctx context.Context, milestone Milestone) error {
	return ErrCollections.Wrap(service.projects.UpdateMilestone(ctx, milestone))
}

// DeleteMilestone deletes certain milestones.
func (service *Service) DeleteMilestone(ctx context.Context, milestoneID uuid.UUID) error {
	return ErrCollections.Wrap(service.projects.DeleteMilestone(ctx, milestoneID))
}

// AddSocialNetwork add social network.
func (service *Service) AddSocialNetwork(ctx context.Context, projectID uuid.UUID, title, link string) error {
	socialNetwork := SocialNetwork{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     title,
		Link:      link,
	}

	return ErrCollections.Wrap(service.projects.AddSocialNetwork(ctx, socialNetwork))
}

// ListSocialNetworks returns collection social networks.
func (service *Service) ListSocialNetworks(ctx context.Context, id uuid.UUID) ([]SocialNetwork, error) {
	socialNetworks, err := service.projects.ListSocialNetworks(ctx, id)

	return socialNetworks, ErrCollections.Wrap(err)
}

// UpdateSocialNetwork updates social network.
func (service *Service) UpdateSocialNetwork(ctx context.Context, socialNetwork SocialNetwork) error {
	return ErrCollections.Wrap(service.projects.UpdateSocialNetwork(ctx, socialNetwork))
}

// DeleteSocialNetwork deletes certain social network.
func (service *Service) DeleteSocialNetwork(ctx context.Context, socialNetworkLinkID uuid.UUID) error {
	return ErrCollections.Wrap(service.projects.DeleteSocialNetwork(ctx, socialNetworkLinkID))
}
