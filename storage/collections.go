package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zeebo/errs"

	"generator/collections"
)

// ensures that collectionsDB implements collections.DB.
var _ collections.DB = (*collectionsDB)(nil)

// ErrCollectionsDB indicates that there was an error in the database.
var ErrCollectionsDB = errs.Class("collections repository error")

// pgUniqueViolation is the PostgreSQL error code of a unique constraint violation.
const pgUniqueViolation = "23505"

// collectionsDB provides access to collections db.
//
// architecture: Database
type collectionsDB struct {
	conn *sql.DB
}

// Create creates collection in database.
func (collectionsDB *collectionsDB) Create(ctx context.Context, project collections.Collection) error {
	query := `INSERT INTO collections(id, name, nft_count, user_id, contract_address, are_images_generated, is_website_published)
	          VALUES($1,$2,$3,$4,$5,$6,$7)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, project.ID, project.Name, project.NFTCount, project.UserID, project.ContractAddress,
		project.AreImagesGenerated, project.IsWebsitePublished)

	return ErrCollectionsDB.Wrap(err)
}

// Get returns collection by id from database.
func (collectionsDB *collectionsDB) Get(ctx context.Context, id uuid.UUID) (collections.Collection, error) {
	var project collections.Collection
	query := `SELECT name, nft_count, user_id, contract_address, contract_block, are_images_generated, is_website_published, COALESCE(website_slug, ''),
	                 images_cid, metadata_cid, ipfs_pinned
	          FROM collections
	          WHERE id=$1`

	err := collectionsDB.conn.QueryRowContext(ctx, query, id).Scan(
		&project.Name, &project.NFTCount, &project.UserID, &project.ContractAddress, &project.ContractBlock, &project.AreImagesGenerated, &project.IsWebsitePublished,
		&project.WebsiteSlug, &project.ImagesCID, &project.MetadataCID, &project.IPFSPinnedIn)
	project.ID = id
	if errs.Is(err, sql.ErrNoRows) {
		return project, collections.ErrNoCollection.Wrap(err)
	}

	return project, ErrCollectionsDB.Wrap(err)
}

// GetByUserID returns all user collections by user id from database.
func (collectionsDB *collectionsDB) GetByUserID(ctx context.Context, id uuid.UUID) (_ []collections.Collection, err error) {
	var userProjects []collections.Collection
	query := `SELECT id, name, nft_count, contract_address, contract_block, are_images_generated, is_website_published, COALESCE(website_slug, ''),
	                 images_cid, metadata_cid, ipfs_pinned
	          FROM collections
	          WHERE user_id=$1`

	rows, err := collectionsDB.conn.QueryContext(ctx, query, id)
	if err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var project collections.Collection
		err = rows.Scan(&project.ID, &project.Name, &project.NFTCount, &project.ContractAddress, &project.ContractBlock, &project.AreImagesGenerated, &project.IsWebsitePublished,
			&project.WebsiteSlug, &project.ImagesCID, &project.MetadataCID, &project.IPFSPinnedIn)
		if err != nil {
			return nil, ErrCollectionsDB.Wrap(err)
		}
		project.UserID = id

		userProjects = append(userProjects, project)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}

	return userProjects, nil
}

// List returns all collections from database.
func (collectionsDB *collectionsDB) List(ctx context.Context) (_ []collections.Collection, err error) {
	var allProjects []collections.Collection
	query := `SELECT id, name, nft_count, user_id, contract_address, contract_block, are_images_generated, is_website_published, COALESCE(website_slug, ''),
	                 images_cid, metadata_cid, ipfs_pinned
	          FROM collections`

	rows, err := collectionsDB.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var project collections.Collection
		err = rows.Scan(&project.ID, &project.Name, &project.NFTCount, &project.UserID, &project.ContractAddress, &project.ContractBlock, &project.AreImagesGenerated,
			&project.IsWebsitePublished, &project.WebsiteSlug, &project.ImagesCID, &project.MetadataCID, &project.IPFSPinnedIn)
		if err != nil {
			return nil, ErrCollectionsDB.Wrap(err)
		}

		allProjects = append(allProjects, project)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}

	return allProjects, nil
}

// Update updates collection in the database.
func (collectionsDB *collectionsDB) Update(ctx context.Context, id uuid.UUID, update collections.Update) error {
	query := `UPDATE collections
	          SET name = $1, nft_count = $2
	          WHERE id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, update.Name, update.NFTCount, id)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoCollection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// UpdateContractAddress updates contract address and its deployment block for the collection in the database.
func (collectionsDB *collectionsDB) UpdateContractAddress(ctx context.Context, id uuid.UUID, contractAddress common.Address, contractBlock uint64) error {
	query := `UPDATE collections
	          SET contract_address = $1, contract_block = $2
	          WHERE id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, contractAddress, int64(contractBlock), id)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoCollection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// GetBySlug returns collection by its website address from database.
func (collectionsDB *collectionsDB) GetBySlug(ctx context.Context, slug string) (collections.Collection, error) {
	project := collections.Collection{WebsiteSlug: slug}
	query := `SELECT id, name, nft_count, user_id, contract_address, contract_block, are_images_generated, is_website_published,
	                 images_cid, metadata_cid, ipfs_pinned
	          FROM collections
	          WHERE website_slug=$1`

	err := collectionsDB.conn.QueryRowContext(ctx, query, slug).Scan(
		&project.ID, &project.Name, &project.NFTCount, &project.UserID, &project.ContractAddress, &project.ContractBlock, &project.AreImagesGenerated,
		&project.IsWebsitePublished, &project.ImagesCID, &project.MetadataCID, &project.IPFSPinnedIn)
	if errs.Is(err, sql.ErrNoRows) {
		return project, collections.ErrNoCollection.Wrap(err)
	}

	return project, ErrCollectionsDB.Wrap(err)
}

// UpdateIPFS updates the IPFS directories of the published collection in the database.
func (collectionsDB *collectionsDB) UpdateIPFS(ctx context.Context, id uuid.UUID, imagesCID, metadataCID string) error {
	// other directories are not pinned anywhere yet.
	query := `UPDATE collections
	          SET ipfs_pinned = CASE WHEN images_cid = $1 AND metadata_cid = $2 THEN ipfs_pinned ELSE '' END,
	              images_cid = $1, metadata_cid = $2
	          WHERE id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, imagesCID, metadataCID, id)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoCollection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// UpdateWebsite updates website address and publication state of the collection in the database.
func (collectionsDB *collectionsDB) UpdateWebsite(ctx context.Context, id uuid.UUID, slug string, isPublished bool) error {
	query := `UPDATE collections
	          SET website_slug = NULLIF($1, ''), is_website_published = $2
	          WHERE id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, slug, isPublished, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgUniqueViolation {
			return collections.ErrSlugTaken.Wrap(err)
		}
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoCollection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// Delete deletes collection from database.
func (collectionsDB *collectionsDB) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM collections
	          WHERE id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, id)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoCollection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// CreateSite creates site for collection in the database.
func (collectionsDB *collectionsDB) CreateSite(ctx context.Context, site collections.Site) error {
	query := `INSERT INTO collection_mint_site(project_id, title_header, description_header, title_about, description_about)
	          VALUES($1,$2,$3,$4,$5)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, site.Header.ProjectID, site.Header.Title, site.Header.Description,
		site.About.Title, site.About.Description)

	return ErrCollectionsDB.Wrap(err)
}

// DeleteSite deletes collection site from the database.
func (collectionsDB *collectionsDB) DeleteSite(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM collection_mint_site
	          WHERE project_id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, id)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoSite.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// GetHeaderSection returns header section of collection site from the database.
func (collectionsDB *collectionsDB) GetHeaderSection(ctx context.Context, id uuid.UUID) (collections.Section, error) {
	var section collections.Section
	query := `SELECT title_header, description_header
	          FROM collection_mint_site
	          WHERE project_id=$1`

	err := collectionsDB.conn.QueryRowContext(ctx, query, id).Scan(&section.Title, &section.Description)
	section.ProjectID = id
	if errs.Is(err, sql.ErrNoRows) {
		return section, collections.ErrNoHeaderSection.Wrap(err)
	}

	return section, ErrCollectionsDB.Wrap(err)
}

// UpdateHeaderSection updates header section of collection site in the database.
func (collectionsDB *collectionsDB) UpdateHeaderSection(ctx context.Context, section collections.Section) error {
	query := `UPDATE collection_mint_site
	          SET title_header = $1, description_header = $2
	          WHERE project_id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, section.Title, section.Description, section.ProjectID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoHeaderSection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// GetAboutSection returns about section of collection site from the database.
func (collectionsDB *collectionsDB) GetAboutSection(ctx context.Context, id uuid.UUID) (collections.Section, error) {
	var section collections.Section
	query := `SELECT title_about, description_about
	          FROM collection_mint_site
	          WHERE project_id=$1`

	err := collectionsDB.conn.QueryRowContext(ctx, query, id).Scan(&section.Title, &section.Description)
	section.ProjectID = id
	if errs.Is(err, sql.ErrNoRows) {
		return section, collections.ErrNoAboutSection.Wrap(err)
	}

	return section, ErrCollectionsDB.Wrap(err)
}

// UpdateAboutSection updates about section of collection site in the database.
func (collectionsDB *collectionsDB) UpdateAboutSection(ctx context.Context, section collections.Section) error {
	query := `UPDATE collection_mint_site
	          SET title_about = $1, description_about = $2
	          WHERE project_id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, section.Title, section.Description, section.ProjectID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoAboutSection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// AddMember adds member of collection in the database.
func (collectionsDB *collectionsDB) AddMember(ctx context.Context, member collections.Member) error {
	query := `INSERT INTO collection_mint_site_members(id, project_id, name, role, created_at)
	          VALUES($1,$2,$3,$4,$5)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, member.ID, member.ProjectID, member.Name, member.Role, time.Now().UTC())

	return ErrCollectionsDB.Wrap(err)
}

// ListMembers returns all collection members from the database.
func (collectionsDB *collectionsDB) ListMembers(ctx context.Context, id uuid.UUID) (_ []collections.Member, err error) {
	var members []collections.Member
	query := `SELECT id, name, role
	          FROM collection_mint_site_members
	          WHERE project_id=$1
	          ORDER BY created_at`

	rows, err := collectionsDB.conn.QueryContext(ctx, query, id)
	if err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var member collections.Member
		err = rows.Scan(&member.ID, &member.Name, &member.Role)
		member.ProjectID = id
		if err != nil {
			return nil, ErrCollectionsDB.Wrap(err)
		}

		members = append(members, member)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}

	return members, nil
}

// GetMember returns collection member from the database.
func (collectionsDB *collectionsDB) GetMember(ctx context.Context, id uuid.UUID) (collections.Member, error) {
	var member collections.Member
	query := `SELECT project_id, name, role
	          FROM collection_mint_site_members
	          WHERE id=$1`

	err := collectionsDB.conn.QueryRowContext(ctx, query, id).Scan(&member.ProjectID, &member.Name, &member.Role)
	if errs.Is(err, sql.ErrNoRows) {
		return member, collections.ErrNoMember.New("")
	}
	member.ID = id

	return member, ErrCollectionsDB.Wrap(err)
}

// UpdateMember updates member of collection in the database.
func (collectionsDB *collectionsDB) UpdateMember(ctx context.Context, member collections.Member) error {
	query := `UPDATE collection_mint_site_members
	          SET name = $1, role = $2
	          WHERE id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, member.Name, member.Role, member.ID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoMember.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// DeleteMember deletes collection member from the database.
func (collectionsDB *collectionsDB) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	query := `DELETE FROM collection_mint_site_members
	          WHERE id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, memberID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoMember.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// DeleteMembers deletes all collection members from the database.
func (collectionsDB *collectionsDB) DeleteMembers(ctx context.Context, projectID uuid.UUID) error {
	query := `DELETE FROM collection_mint_site_members
	          WHERE project_id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, projectID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoMember.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// AddMilestone adds milestone of collection in the database.
func (collectionsDB *collectionsDB) AddMilestone(ctx context.Context, milestone collections.Milestone) error {
	query := `INSERT INTO collection_mint_site_milestones(id, project_id, title, date, description)
	          VALUES($1,$2,$3,$4,$5)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, milestone.ID, milestone.ProjectID, milestone.Title, milestone.Date, milestone.Description)

	return ErrCollectionsDB.Wrap(err)
}

// ListMilestones returns all collection milestones from the database.
func (collectionsDB *collectionsDB) ListMilestones(ctx context.Context, id uuid.UUID) (_ []collections.Milestone, err error) {
	var milestones []collections.Milestone
	query := `SELECT id, title, date, description
	          FROM collection_mint_site_milestones
	          WHERE project_id=$1`

	rows, err := collectionsDB.conn.QueryContext(ctx, query, id)
	if err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var milestone collections.Milestone
		err = rows.Scan(&milestone.ID, &milestone.Title,
			&milestone.Date, &milestone.Description)
		milestone.ProjectID = id
		if err != nil {
			return nil, ErrCollectionsDB.Wrap(err)
		}

		milestones = append(milestones, milestone)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}

	return milestones, nil
}

// UpdateMilestone updates milestone of collection in the database.
func (collectionsDB *collectionsDB) UpdateMilestone(ctx context.Context, milestone collections.Milestone) error {
	query := `UPDATE collection_mint_site_milestones
	          SET title = $1, date = $2, description = $3
	          WHERE id = $4`

	result, err := collectionsDB.conn.ExecContext(ctx, query, milestone.Title, milestone.Date,
		milestone.Description, milestone.ID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoMilestone.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// DeleteMilestone deletes collection milestone from the database.
func (collectionsDB *collectionsDB) DeleteMilestone(ctx context.Context, milestoneID uuid.UUID) error {
	query := `DELETE FROM collection_mint_site_milestones
	          WHERE id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, milestoneID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoMilestone.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// DeleteMilestones deletes all collection milestones.
func (collectionsDB *collectionsDB) DeleteMilestones(ctx context.Context, projectID uuid.UUID) error {
	query := `DELETE FROM collection_mint_site_milestones
	          WHERE project_id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, projectID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoMilestone.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// AddSocialNetwork adds social networks of collection site in the database.
func (collectionsDB *collectionsDB) AddSocialNetwork(ctx context.Context, socialNetwork collections.SocialNetwork) error {
	query := `INSERT INTO collection_mint_site_social_networks(id, project_id, title, link)
	          VALUES($1,$2,$3,$4)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, socialNetwork.ID, socialNetwork.ProjectID, socialNetwork.Title, socialNetwork.Link)

	return ErrCollectionsDB.Wrap(err)
}

// ListSocialNetworks returns social networks section of collection site from the database.
func (collectionsDB *collectionsDB) ListSocialNetworks(ctx context.Context, id uuid.UUID) (_ []collections.SocialNetwork, err error) {
	var socialNetworks []collections.SocialNetwork
	query := `SELECT id, title, link
	          FROM collection_mint_site_social_networks
	          WHERE project_id=$1`

	rows, err := collectionsDB.conn.QueryContext(ctx, query, id)
	if err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var socialNetwork collections.SocialNetwork
		err = rows.Scan(&socialNetwork.ID, &socialNetwork.Title, &socialNetwork.Link)
		socialNetwork.ProjectID = id
		if err != nil {
			return nil, ErrCollectionsDB.Wrap(err)
		}

		socialNetworks = append(socialNetworks, socialNetwork)
	}
	if err = rows.Err(); err != nil {
		return nil, ErrCollectionsDB.Wrap(err)
	}

	return socialNetworks, nil
}

// UpdateSocialNetwork updates social networks of collection site in the database.
func (collectionsDB *collectionsDB) UpdateSocialNetwork(ctx context.Context, socialNetwork collections.SocialNetwork) error {
	query := `UPDATE collection_mint_site_social_networks
	          SET title = $1, link = $2
	          WHERE id = $3`

	result, err := collectionsDB.conn.ExecContext(ctx, query, socialNetwork.Title,
		socialNetwork.Link, socialNetwork.ID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoSocialNetwork.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// DeleteSocialNetwork deletes collection social network from the database.
func (collectionsDB *collectionsDB) DeleteSocialNetwork(ctx context.Context, socialNetworkID uuid.UUID) error {
	query := `DELETE FROM collection_mint_site_social_networks
	          WHERE id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, socialNetworkID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoSocialNetwork.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// DeleteSocialNetworks deletes all collection social networks from the database.
func (collectionsDB *collectionsDB) DeleteSocialNetworks(ctx context.Context, projectID uuid.UUID) error {
	query := `DELETE FROM collection_mint_site_social_networks
	          WHERE project_id = $1`

	result, err := collectionsDB.conn.ExecContext(ctx, query, projectID)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoSocialNetwork.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}

// UpdateIPFSPinned saves the pinning service that keeps the published directories of the collection.
func (collectionsDB *collectionsDB) UpdateIPFSPinned(ctx context.Context, id uuid.UUID, service string) error {
	query := `UPDATE collections
	          SET ipfs_pinned = $1
	          WHERE id = $2`

	result, err := collectionsDB.conn.ExecContext(ctx, query, service, id)
	if err != nil {
		return ErrCollectionsDB.Wrap(err)
	}
	rowNum, err := result.RowsAffected()
	if err == nil && rowNum == 0 {
		return collections.ErrNoCollection.New("")
	}

	return ErrCollectionsDB.Wrap(err)
}
