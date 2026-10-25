package collections_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator"
	"generator/accounts"
	"generator/collections"
	"generator/storage/testdb"
)

func TestProjects(t *testing.T) {
	user1 := accounts.Account{
		ID:           uuid.New(),
		Email:        "testUser1@gmail.com",
		PasswordHash: []byte("Qwerty123-"),
		LastLogin:    time.Now().UTC(),
		Status:       accounts.StatusCreated,
		CreatedAt:    time.Now().UTC(),
	}

	project1 := collections.Collection{
		ID:                 uuid.New(),
		Name:               "project1",
		NFTCount:           10,
		UserID:             user1.ID,
		ContractAddress:    common.HexToAddress("0x"),
		AreImagesGenerated: false,
		IsWebsitePublished: false,
	}

	project2 := collections.Collection{
		ID:                 uuid.New(),
		Name:               "project2",
		NFTCount:           10,
		UserID:             user1.ID,
		ContractAddress:    common.HexToAddress("0x"),
		AreImagesGenerated: false,
		IsWebsitePublished: false,
	}

	updatedProject1 := collections.Update{
		Name:     "updatedProject1",
		NFTCount: 100,
	}

	aboutSection := collections.Section{
		ProjectID:   project1.ID,
		Title:       "test about section",
		Description: "test about description",
	}

	headerSection := collections.Section{
		ProjectID:   project1.ID,
		Title:       "test header section",
		Description: "test header description",
	}

	site := collections.Site{
		ProjectID: project1.ID,
		About:     aboutSection,
		Header:    headerSection,
	}

	member1 := collections.Member{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "member1",
		Role:      "test",
	}

	member2 := collections.Member{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Name:      "member2",
		Role:      "test",
	}

	milestone1 := collections.Milestone{
		ID:          uuid.New(),
		ProjectID:   project1.ID,
		Title:       "test1",
		Date:        "date1",
		Description: "test1",
	}

	milestone2 := collections.Milestone{
		ID:          uuid.New(),
		ProjectID:   project1.ID,
		Title:       "test2",
		Date:        "date2",
		Description: "test2",
	}

	socialNetwork1 := collections.SocialNetwork{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Title:     "twitter",
		Link:      "twitter.com",
	}

	socialNetwork2 := collections.SocialNetwork{
		ID:        uuid.New(),
		ProjectID: project1.ID,
		Title:     "instagram",
		Link:      "instagram.com",
	}

	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		repositoryUsers := db.Accounts()
		repositoryProjects := db.Collections()

		t.Run("create", func(t *testing.T) {
			err := repositoryUsers.Create(ctx, user1)
			require.NoError(t, err)

			err = repositoryProjects.Create(ctx, project1)
			require.NoError(t, err)

			err = repositoryProjects.Create(ctx, project2)
			require.NoError(t, err)
		})

		t.Run("get sql no rows", func(t *testing.T) {
			_, err := repositoryProjects.Get(ctx, uuid.New())
			require.Error(t, err)
			assert.True(t, collections.ErrNoCollection.Has(err))
		})

		t.Run("get", func(t *testing.T) {
			project, err := repositoryProjects.Get(ctx, project1.ID)
			require.NoError(t, err)
			compareProjects(t, project1, project)
		})

		t.Run("list", func(t *testing.T) {
			allProjects, err := repositoryProjects.List(ctx)
			require.NoError(t, err)
			require.Equal(t, len(allProjects), 2)
			compareProjects(t, allProjects[0], project1)
			compareProjects(t, allProjects[1], project2)
		})

		t.Run("get user projects", func(t *testing.T) {
			userProjects, err := repositoryProjects.GetByUserID(ctx, user1.ID)
			require.NoError(t, err)
			require.Equal(t, len(userProjects), 2)
			compareProjects(t, userProjects[0], project1)
			compareProjects(t, userProjects[1], project2)
		})

		t.Run("update sql no rows", func(t *testing.T) {
			err := repositoryProjects.Update(ctx, uuid.New(), collections.Update{})
			require.Error(t, err)
			assert.True(t, collections.ErrNoCollection.Has(err))
		})

		t.Run("update", func(t *testing.T) {
			err := repositoryProjects.Update(ctx, project1.ID, updatedProject1)
			require.NoError(t, err)

			project, err := repositoryProjects.Get(ctx, project1.ID)
			require.NoError(t, err)
			project1.Name = updatedProject1.Name
			project1.NFTCount = updatedProject1.NFTCount
			compareProjects(t, project, project1)
		})

		t.Run("UpdateContractAddress", func(t *testing.T) {
			project1.ContractAddress = common.HexToAddress("0x2df71a8114129312c0e5e26c26a9bc61e22d6fd7")
			project1.ContractBlock = 11891261
			err := repositoryProjects.UpdateContractAddress(ctx, project1.ID, project1.ContractAddress, project1.ContractBlock)
			require.NoError(t, err)

			project, err := repositoryProjects.Get(ctx, project1.ID)
			require.NoError(t, err)
			compareProjects(t, project, project1)
		})

		t.Run("UpdateWebsite and GetBySlug", func(t *testing.T) {
			_, err := repositoryProjects.GetBySlug(ctx, "project-one")
			require.Error(t, err)
			assert.True(t, collections.ErrNoCollection.Has(err))

			err = repositoryProjects.UpdateWebsite(ctx, uuid.New(), "project-one", true)
			assert.True(t, collections.ErrNoCollection.Has(err))

			err = repositoryProjects.UpdateWebsite(ctx, project1.ID, "project-one", true)
			require.NoError(t, err)
			project1.WebsiteSlug, project1.IsWebsitePublished = "project-one", true

			project, err := repositoryProjects.GetBySlug(ctx, "project-one")
			require.NoError(t, err)
			compareProjects(t, project, project1)
			assert.Equal(t, "project-one", project.WebsiteSlug)

			project, err = repositoryProjects.Get(ctx, project1.ID)
			require.NoError(t, err)
			assert.Equal(t, "project-one", project.WebsiteSlug)

			// The address belongs to project1, so project2 cannot take it.
			err = repositoryProjects.UpdateWebsite(ctx, project2.ID, "project-one", true)
			require.Error(t, err)
			assert.True(t, collections.ErrSlugTaken.Has(err))

			// Collections without an address do not conflict with each other.
			err = repositoryProjects.UpdateWebsite(ctx, project2.ID, "", false)
			require.NoError(t, err)

			err = repositoryProjects.UpdateWebsite(ctx, project1.ID, "project-one", false)
			require.NoError(t, err)
			project1.IsWebsitePublished = false
		})

		t.Run("delete sql no rows", func(t *testing.T) {
			err := repositoryProjects.Delete(ctx, uuid.New())
			require.Error(t, err)
			assert.True(t, collections.ErrNoCollection.Has(err))
		})

		t.Run("delete", func(t *testing.T) {
			err := repositoryProjects.Delete(ctx, project1.ID)
			require.NoError(t, err)

			err = repositoryProjects.Delete(ctx, project2.ID)
			require.NoError(t, err)
		})

		t.Run("create site", func(t *testing.T) {
			err := repositoryProjects.Create(ctx, project1)
			require.NoError(t, err)

			err = repositoryProjects.CreateSite(ctx, site)
			require.NoError(t, err)
		})

		t.Run("unsuccessful delete site", func(t *testing.T) {
			err := repositoryProjects.DeleteSite(ctx, uuid.Nil)
			require.Error(t, err)
			assert.True(t, collections.ErrNoSite.Has(err))
		})

		t.Run("delete site", func(t *testing.T) {
			err := repositoryProjects.DeleteSite(ctx, site.ProjectID)
			require.NoError(t, err)

			_, err = repositoryProjects.GetHeaderSection(ctx, site.ProjectID)
			require.Error(t, err)
			assert.True(t, collections.ErrNoHeaderSection.Has(err))
		})

		t.Run("update header section sql no rows", func(t *testing.T) {
			err := repositoryProjects.CreateSite(ctx, site)
			require.NoError(t, err)

			err = repositoryProjects.UpdateHeaderSection(ctx, collections.Section{})
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoHeaderSection.Has(err))
		})

		t.Run("update header section", func(t *testing.T) {
			err := repositoryProjects.UpdateHeaderSection(ctx, headerSection)
			require.NoError(t, err)
		})

		t.Run("get header section sql no rows", func(t *testing.T) {
			_, err := repositoryProjects.GetHeaderSection(ctx, uuid.New())
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoHeaderSection.Has(err))
		})

		t.Run("get header section", func(t *testing.T) {
			headerSectionDB, err := repositoryProjects.GetHeaderSection(ctx, headerSection.ProjectID)
			require.NoError(t, err)

			compareSection(t, headerSectionDB, headerSection)
		})

		t.Run("update about section sql no rows", func(t *testing.T) {
			err := repositoryProjects.UpdateAboutSection(ctx, collections.Section{})
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoAboutSection.Has(err))
		})

		t.Run("update about section", func(t *testing.T) {
			err := repositoryProjects.UpdateAboutSection(ctx, aboutSection)
			require.NoError(t, err)
		})

		t.Run("get about section sql no rows", func(t *testing.T) {
			_, err := repositoryProjects.GetAboutSection(ctx, uuid.New())
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoAboutSection.Has(err))
		})

		t.Run("get about section", func(t *testing.T) {
			aboutSectionDB, err := repositoryProjects.GetAboutSection(ctx, aboutSection.ProjectID)
			require.NoError(t, err)
			compareSection(t, aboutSectionDB, aboutSection)
		})

		t.Run("create member", func(t *testing.T) {
			err := repositoryProjects.AddMember(ctx, member1)
			require.NoError(t, err)
		})

		t.Run("list member", func(t *testing.T) {
			err := repositoryProjects.AddMember(ctx, member2)
			require.NoError(t, err)

			members, err := repositoryProjects.ListMembers(ctx, member1.ProjectID)
			require.NoError(t, err)
			require.Equal(t, 2, len(members))
			compareMembers(t, members[0], member1)
			compareMembers(t, members[1], member2)
		})

		t.Run("negative list members", func(t *testing.T) {
			members, err := repositoryProjects.ListMembers(ctx, member1.ProjectID)
			require.NoError(t, err)

			// compares members from db with invalid member order.
			// comparing names and roles unnecessary.
			require.Equal(t, 2, len(members))
			assert.False(t, members[0].ID == member2.ID)
			assert.False(t, members[1].ID == member1.ID)
		})

		t.Run("update member sql no rows", func(t *testing.T) {
			err := repositoryProjects.UpdateMember(ctx, collections.Member{})
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoMember.Has(err))
		})

		t.Run("update member", func(t *testing.T) {
			member1.Name = "update"
			member1.Role = "updated"

			err := repositoryProjects.UpdateMember(ctx, member1)
			require.NoError(t, err)
		})

		t.Run("delete member sql no rows", func(t *testing.T) {
			err := repositoryProjects.DeleteMember(ctx, uuid.New())
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoMember.Has(err))
		})

		t.Run("delete member", func(t *testing.T) {
			err := repositoryProjects.DeleteMember(ctx, member1.ID)
			require.NoError(t, err)
		})

		t.Run("unsuccessful delete member", func(t *testing.T) {
			err := repositoryProjects.DeleteMembers(ctx, uuid.Nil)
			require.Error(t, err)
			require.True(t, collections.ErrNoMember.Has(err))
		})

		t.Run("delete members", func(t *testing.T) {
			err := repositoryProjects.DeleteMembers(ctx, member1.ProjectID)
			require.NoError(t, err)
		})

		t.Run("create milestone", func(t *testing.T) {
			err := repositoryProjects.AddMilestone(ctx, milestone1)
			require.NoError(t, err)
		})

		t.Run("list milestone", func(t *testing.T) {
			err := repositoryProjects.AddMilestone(ctx, milestone2)
			require.NoError(t, err)

			milestones, err := repositoryProjects.ListMilestones(ctx, milestone2.ProjectID)
			require.NoError(t, err)

			require.Equal(t, 2, len(milestones))
			compareMilestone(t, milestones[0], milestone1)
			compareMilestone(t, milestones[1], milestone2)
		})

		t.Run("update milestone sql no rows", func(t *testing.T) {
			err := repositoryProjects.UpdateMilestone(ctx, collections.Milestone{})
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoMilestone.Has(err))
		})

		t.Run("update milestone", func(t *testing.T) {
			milestone1.Title = "updated"
			milestone1.Description = "updated"
			milestone1.Date = "test3"

			err := repositoryProjects.UpdateMilestone(ctx, milestone1)
			require.NoError(t, err)
		})

		t.Run("delete milestone sql no rows", func(t *testing.T) {
			err := repositoryProjects.DeleteMilestone(ctx, uuid.New())
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoMilestone.Has(err))
		})

		t.Run("delete milestone", func(t *testing.T) {
			err := repositoryProjects.DeleteMilestone(ctx, milestone1.ID)
			require.NoError(t, err)
		})

		t.Run("unsuccessful delete milestones", func(t *testing.T) {
			err := repositoryProjects.DeleteMilestone(ctx, uuid.Nil)
			require.Error(t, err)
			assert.True(t, collections.ErrNoMilestone.Has(err))
		})

		t.Run("delete milestones", func(t *testing.T) {
			err := repositoryProjects.DeleteMilestones(ctx, milestone1.ProjectID)
			require.NoError(t, err)
		})

		t.Run("create social network", func(t *testing.T) {
			err := repositoryProjects.AddSocialNetwork(ctx, socialNetwork1)
			require.NoError(t, err)
		})

		t.Run("list social network", func(t *testing.T) {
			err := repositoryProjects.AddSocialNetwork(ctx, socialNetwork2)
			require.NoError(t, err)

			socialNetworks, err := repositoryProjects.ListSocialNetworks(ctx, socialNetwork1.ProjectID)
			require.NoError(t, err)
			require.Equal(t, 2, len(socialNetworks))
			compareSocialNetworks(t, socialNetworks[0], socialNetwork1)
			compareSocialNetworks(t, socialNetworks[1], socialNetwork2)
		})

		t.Run("update social network sql no rows", func(t *testing.T) {
			err := repositoryProjects.UpdateSocialNetwork(ctx, collections.SocialNetwork{})
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoSocialNetwork.Has(err))
		})

		t.Run("update social network", func(t *testing.T) {
			socialNetwork1.Title = "updated"
			socialNetwork1.Link = "updated.com"

			err := repositoryProjects.UpdateSocialNetwork(ctx, socialNetwork1)
			require.NoError(t, err)
		})

		t.Run("delete social network sql no rows", func(t *testing.T) {
			err := repositoryProjects.DeleteSocialNetwork(ctx, uuid.New())
			require.Error(t, err)
			assert.Equal(t, true, collections.ErrNoSocialNetwork.Has(err))
		})

		t.Run("delete social network", func(t *testing.T) {
			err := repositoryProjects.DeleteSocialNetwork(ctx, socialNetwork1.ID)
			require.NoError(t, err)
		})

		t.Run("unsuccessful delete social networks", func(t *testing.T) {
			err := repositoryProjects.DeleteSocialNetworks(ctx, uuid.Nil)
			require.Error(t, err)
			assert.True(t, collections.ErrNoSocialNetwork.Has(err))
		})

		t.Run("delete social networks", func(t *testing.T) {
			err := repositoryProjects.DeleteSocialNetworks(ctx, socialNetwork2.ProjectID)
			require.NoError(t, err)
		})
	})
}

func compareProjects(t *testing.T, project, projectToCompare collections.Collection) {
	assert.Equal(t, project.ID, projectToCompare.ID)
	assert.Equal(t, project.Name, projectToCompare.Name)
	assert.Equal(t, project.NFTCount, projectToCompare.NFTCount)
	assert.Equal(t, project.UserID, projectToCompare.UserID)
	assert.Equal(t, project.ContractAddress, projectToCompare.ContractAddress)
	assert.Equal(t, project.ContractBlock, projectToCompare.ContractBlock)
	assert.Equal(t, project.AreImagesGenerated, projectToCompare.AreImagesGenerated)
	assert.Equal(t, project.IsWebsitePublished, projectToCompare.IsWebsitePublished)
}

func compareSection(t *testing.T, section, sectionToCompare collections.Section) {
	assert.Equal(t, section.ProjectID, sectionToCompare.ProjectID)
	assert.Equal(t, section.Title, sectionToCompare.Title)
	assert.Equal(t, section.Description, sectionToCompare.Description)
}

func compareMembers(t *testing.T, member, memberToComparer collections.Member) {
	assert.Equal(t, member.ID, memberToComparer.ID)
	assert.Equal(t, member.ProjectID, memberToComparer.ProjectID)
	assert.Equal(t, member.Name, memberToComparer.Name)
	assert.Equal(t, member.Role, memberToComparer.Role)
}

func compareMilestone(t *testing.T, milestone, milestoneToComparer collections.Milestone) {
	assert.Equal(t, milestone.ID, milestoneToComparer.ID)
	assert.Equal(t, milestone.ProjectID, milestoneToComparer.ProjectID)
	assert.Equal(t, milestone.Title, milestoneToComparer.Title)
	assert.Equal(t, milestone.Description, milestoneToComparer.Description)
	assert.Equal(t, milestone.Date, milestoneToComparer.Date)
}

func compareSocialNetworks(t *testing.T, socialNetwork, socialNetworkToCompare collections.SocialNetwork) {
	assert.Equal(t, socialNetwork.ID, socialNetworkToCompare.ID)
	assert.Equal(t, socialNetwork.ProjectID, socialNetworkToCompare.ProjectID)
	assert.Equal(t, socialNetwork.Title, socialNetworkToCompare.Title)
	assert.Equal(t, socialNetwork.Link, socialNetworkToCompare.Link)
}
