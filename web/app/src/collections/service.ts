import { ProjectsClient } from '@/api/projects';
import { AboutSection, CreateMilestoneFields, Header, Member, Milestone, Mintsite, Project, SocialNetwork } from '@/collections';

/**
 * Exposes all projects domain entity related logic.
 */
export class ProjectsService {
    protected readonly projectsClient: ProjectsClient;

    /** Contains http implementation of projects API. */
    public constructor(projectsClient: ProjectsClient) {
        this.projectsClient = projectsClient;
    };

    /** Creates new project. */
    public async create(name: string, nftCount: number): Promise<void> {
        await this.projectsClient.create(name, nftCount);
    };

    /** Deletes project. */
    public async delete(id: string): Promise<void> {
        await this.projectsClient.delete(id);
    };

    /** Requests project. */
    public async getProject(id: string): Promise<Project> {
        return await this.projectsClient.getProject(id);
    };

    /** Requests list of projects. */
    public async list(): Promise<Project[]> {
        return await this.projectsClient.list();
    };

    /** Saves the deployed contract address and the number of the block it was deployed in. */
    public async updateContractAddress(projectId: string, address: string, block: number): Promise<void> {
        await this.projectsClient.updateContractAddress(projectId, address, block);
    }

    /** Updates project by id. */
    public async updateProject(id: string, project: Project): Promise<void> {
        await this.projectsClient.updateProject(id, project);
    };

    /** Generates nfts for current project */
    public async generateNfts(id: string): Promise<void> {
        await this.projectsClient.generateNfts(id);
    };

    /** Deletes generated NTF's. */
    public async deleteGeneratedNFTs(id: string): Promise<void> {
        await this.projectsClient.deleteGeneratedNFTs(id);
    };

    /** Aborts NFT's generation. */
    public async abortGenerateNFTs(id: string): Promise<void> {
        await this.projectsClient.abortGenerateNFTs(id);
    };

    /** Returns progress of current nft generation */
    public async getGenerationProgress(id: string): Promise<number> {
        return await this.projectsClient.getGenerationProgress(id);
    }

    /** Generates metadata for the generated collection. */
    public async generateCollectionMetadata(id: string): Promise<void> {
        await this.projectsClient.generateCollectionMetadata(id);
    }

    /** Publishes the generated collection to IPFS. */
    public async publishToIPFS(id: string): Promise<Project> {
        return await this.projectsClient.publishToIPFS(id);
    }

    /** Pins the published collection in the configured pinning service. */
    public async pinIPFS(id: string): Promise<Project> {
        return await this.projectsClient.pinIPFS(id);
    }

    /** Publishes project mint site at /sites/{slug}. */
    public async publishWebsite(projectId: string, slug: string): Promise<void> {
        await this.projectsClient.publishWebsite(projectId, slug);
    };

    /** Requests project mintsite. */
    public async getProjectMintSite(projectId: string): Promise<Mintsite> {
        return await this.projectsClient.getProjectMintSite(projectId);
    };

    /** Deletes project mintsite. */
    public async deleteMintSite(projectId: string): Promise<void> {
        return await this.projectsClient.deleteMintSite(projectId);
    };

    /** Creates project site 'header' section. */
    public async createHeader(projectId: string, title: string, description: string): Promise<void> {
        await this.projectsClient.createHeader(projectId, title, description);
    };

    /** Request 'header' section content. */
    public async getHeader(projectId: string): Promise<Header> {
        return await this.projectsClient.getHeader(projectId);
    };

    /** Updates 'header' logo. */
    public async updateHeaderLogo(header: Header): Promise<void> {
        return await this.projectsClient.updateHeaderLogo(header);
    };

    /** Updates 'header' image. */
    public async updateHeaderImage(header: Header): Promise<void> {
        return await this.projectsClient.updateHeaderImage(header);
    };

    /** Updates 'header' text fields. */
    public async updateHeaderTextFields(header: Header): Promise<void> {
        return await this.projectsClient.updateHeaderTextFields(header);
    };

    /** Request 'about' section content. */
    public async getAboutSection(projectId: string): Promise<AboutSection> {
        return await this.projectsClient.getAboutSection(projectId);
    };

    /** Creates project site 'about' section. */
    public async createAboutSection(projectId: string, title: string, description: string): Promise<void> {
        await this.projectsClient.createAboutSection(projectId, title, description);
    };

    /** Updates project site 'about' section. */
    public async updateAboutSection(aboutSection: AboutSection): Promise<void> {
        await this.projectsClient.updateAboutSection(aboutSection);
    };

    /** Updates project site 'about' section image. */
    public async updateAboutSectionImage(projectId: string, image: string): Promise<void> {
        await this.projectsClient.updateAboutSectionImage(projectId, image);
    };

    /** Creates new member. */
    public async createMember(projectId: string, name: string, role: string): Promise<void> {
        await this.projectsClient.createMember(projectId, name, role);
    };

    /** Returns members list. */
    public async getMembers(projectId: string): Promise<Member[]> {
        return await this.projectsClient.getMembers(projectId);
    };

    /** Updates member. */
    public async updateMember(member: Member): Promise<void> {
        await this.projectsClient.updateMember(member);
    };

    /** Updates member image. */
    public async updateMemberImage(projectId: string, memberId: string, image: string): Promise<void> {
        await this.projectsClient.updateMemberImage(projectId, memberId, image);
    };

    /** Deletes member. */
    public async deleteMember(projectId: string, memberId: string): Promise<void> {
        await this.projectsClient.deleteMember(projectId, memberId);
    };

    /** Deletes member image. */
    public async deleteMemberImage(projectId: string, memberId: string): Promise<void> {
        await this.projectsClient.deleteMemberImage(projectId, memberId);
    };

    /** Creates new milestone. */
    public async createMilestone(projectId: string, milestoneFields: CreateMilestoneFields): Promise<void> {
        await this.projectsClient.createMilestone(projectId, milestoneFields);
    };

    /** Returns milestones list. */
    public async getMilestones(projectId: string): Promise<Milestone[]> {
        return await this.projectsClient.getMilestones(projectId);
    };

    /** Updates milestone. */
    public async updateMilestone(milestone: Milestone): Promise<void> {
        await this.projectsClient.updateMilestone(milestone);
    };

    /** Deletes milestone. */
    public async deleteMilestone(projectId: string, milestoneId: string): Promise<void> {
        await this.projectsClient.deleteMilestone(projectId, milestoneId);
    };
    /** Creates new social network. */
    public async createSocialNetwork(projectId: string, title: string, link: string): Promise<void> {
        await this.projectsClient.createSocialNetwork(projectId, title, link);
    };

    /** Returns social networks list. */
    public async getSocialNetworks(projectId: string): Promise<SocialNetwork[]> {
        return await this.projectsClient.getSocialNetworks(projectId);
    };

    /** Updates social network. */
    public async updateSocialNetwork(socialNetwork: SocialNetwork): Promise<void> {
        await this.projectsClient.updateSocialNetwork(socialNetwork);
    };

    /** Deletes social network. */
    public async deleteSocialNetwork(projectId: string, socialNetworkId: string): Promise<void> {
        await this.projectsClient.deleteSocialNetwork(projectId, socialNetworkId);
    };

    /** Updates project site 'content' image. */
    public async updateContentImage(projectId: string, image: string, name: string): Promise<void> {
        await this.projectsClient.updateContentImage(projectId, image, name);
    };

    /** Gets project site 'content' images. */
    public async getContentImageNames(projectId: string): Promise<string[]> {
        return await this.projectsClient.getContentImageNames(projectId);
    };

    /** Deletes project site 'content' image. */
    public async deleteContentImage(projectId: string, name: string): Promise<void> {
        await this.projectsClient.deleteContentImage(projectId, name);
    };
};
