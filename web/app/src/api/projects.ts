import { APIClient } from '@/api';
import { AboutSection, CreateMilestoneFields, Header, Member, Milestone, Mintsite, Project, SocialNetwork } from '@/collections';

// TODO: also this logic could be changed after back-end API solution.
/**
 * ProjectsClient is a http implementation of projects API.
 * Exposes all projects-related functionality.
 */
export class ProjectsClient extends APIClient {
    private readonly ROOT_PATH: string = '/api/v0/collections';

    /** Creates new project. */
    public async create(name: string, nftCount: number): Promise<void> {
        const path = `${this.ROOT_PATH}`;

        const response = await this.http.post(path, JSON.stringify({ name, nftCount }));

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Deletes project. */
    public async delete(id: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${id}`;

        const response = await this.http.delete(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Requests project. */
    public async getProject(id: string): Promise<Project> {
        const path = `${this.ROOT_PATH}/${id}`;

        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const project = await response.json();

        return new Project(
            project.id,
            project.userId,
            project.name,
            project.contractAddress,
            project.nftCount,
            project.isWebsitePublished,
            project.areImagesGenerated,
            project.websiteSlug,
            project.imagesCid,
            project.metadataCid,
            project.ipfsPinnedIn,
        );
    };

    /** Requests a list of projects. */
    public async list(): Promise<Project[]> {
        const path = `${this.ROOT_PATH}`;

        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const projects = await response.json() || [];

        return projects.map((project: Project) => new Project(
            project.id,
            project.userId,
            project.name,
            project.contractAddress,
            project.nftCount,
            project.isWebsitePublished,
            project.areImagesGenerated,
            project.websiteSlug,
            project.imagesCid,
            project.metadataCid,
            project.ipfsPinnedIn,
        ));
    };

    /** Saves the deployed contract address and the number of the block it was deployed in. */
    public async updateContractAddress(projectId: string, contractAddress: string, contractBlock: number): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/contractAddress`;

        const response = await this.http.patch(path, JSON.stringify({ contractAddress, contractBlock }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Publishes project mint site at /sites/{slug}. */
    public async publishWebsite(projectId: string, slug: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/website`;

        const response = await this.http.put(path, JSON.stringify({ slug }));

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Requests project mintsite. */
    public async getProjectMintSite(projectId: string): Promise<Mintsite> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site`;

        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const mintsite = await response.json();

        return new Mintsite(mintsite.aboutSelection, mintsite.team, mintsite.roadmap, mintsite.socialNetwork);
    };

    /** Deletes project mintsite. */
    public async deleteMintSite(projectId: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site`;
        const response = await this.http.delete(path);
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Updates project. */
    public async updateProject(id: string, project: Project): Promise<void> {
        const path = `${this.ROOT_PATH}/${id}`;

        const response = await this.http.patch(path, JSON.stringify(project));

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Generates nfts for current project */
    public async generateNfts(id: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${id}/generate-images`;
        const response = await this.http.post(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Deletes generated NTF's. */
    public async deleteGeneratedNFTs(id: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${id}/nfts`;
        const response = await this.http.delete(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Returns progress of current nft generation */
    public async getGenerationProgress(id: string): Promise<number> {
        const path = `${this.ROOT_PATH}/${id}/generated-images-count`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        return await response.json();
    };


    /** Aborts NFT's generation. */
    public async abortGenerateNFTs(id: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${id}/abort-generation`;
        const response = await this.http.put(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Generates metadata for the generated collection. */
    public async generateCollectionMetadata(id: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${id}/generate-metadata`;
        const response = await this.http.post(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Publishes the generated collection to IPFS and returns the project with the images and metadata CIDs. */
    public async publishToIPFS(id: string): Promise<Project> {
        const path = `${this.ROOT_PATH}/${id}/ipfs`;
        const response = await this.http.post(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const project = await response.json();

        return new Project(
            project.id,
            project.userId,
            project.name,
            project.contractAddress,
            project.nftCount,
            project.isWebsitePublished,
            project.areImagesGenerated,
            project.websiteSlug,
            project.imagesCid,
            project.metadataCid,
            project.ipfsPinnedIn,
        );
    };

    /** Pins the published IPFS directories of the project in the configured pinning service. */
    public async pinIPFS(id: string): Promise<Project> {
        const path = `${this.ROOT_PATH}/${id}/ipfs/pin`;
        const response = await this.http.post(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const project = await response.json();

        return new Project(
            project.id,
            project.userId,
            project.name,
            project.contractAddress,
            project.nftCount,
            project.isWebsitePublished,
            project.areImagesGenerated,
            project.websiteSlug,
            project.imagesCid,
            project.metadataCid,
            project.ipfsPinnedIn,
        );
    };

    /** Creates project site 'header' section. */
    public async createHeader(projectId: string, title: string, description: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/header-section`;
        const response = await this.http.patch(path, JSON.stringify({
            title,
            description,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Request 'header' section content. */
    public async getHeader(projectId: string): Promise<Header> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/header-section`;
        const response = await this.http.get(path);
        if (!response.ok) {
            await this.handleError(response);
        }
        const header = await response.json();

        return new Header(header.projectId, header.title, header.description);
    };

    /** Updates project site 'header' text fields. */
    public async updateHeaderTextFields(header: Header): Promise<void> {
        const path = `${this.ROOT_PATH}/${header.projectId}/mint-site/header-section`;
        const response = await this.http.patch(path, JSON.stringify({
            title: header.title,
            description: header.description,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Updates project site 'header' logo. */
    public async updateHeaderLogo(header: Header): Promise<void> {
        const path = `${this.ROOT_PATH}/${header.projectId}/mint-site/logo`;
        const response = await this.http.put(path, JSON.stringify({
            image: header.logo,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Updates project site 'header' image. */
    public async updateHeaderImage(header: Header): Promise<void> {
        const path = `${this.ROOT_PATH}/${header.projectId}/mint-site/header-section`;
        const response = await this.http.put(path, JSON.stringify({
            image: header.image,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Request 'about' section content. */
    public async getAboutSection(projectId: string): Promise<AboutSection> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/about-section`;
        const response = await this.http.get(path);
        if (!response.ok) {
            await this.handleError(response);
        }
        const aboutSection = await response.json();

        return new AboutSection(aboutSection.projectId, aboutSection.title, aboutSection.description);
    };

    /** Creates project site 'about' section. */
    public async createAboutSection(projectId: string, title: string, description: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/about-section`;
        const response = await this.http.patch(path, JSON.stringify({
            title,
            description,
        }));

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Updates project site 'about' section. */
    public async updateAboutSection(aboutSection: AboutSection) {
        const path = `${this.ROOT_PATH}/${aboutSection.projectId}/mint-site/about-section`;
        const response = await this.http.patch(path, JSON.stringify({
            title: aboutSection.title,
            description: aboutSection.description,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Updates project site 'about' section image. */
    public async updateAboutSectionImage(projectId: string, image: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/about-section`;
        // TODO: string is valid json. It will be changed to JSON.stringify(image).
        const response = await this.http.put(path, JSON.stringify({ image }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Creates new member. */
    public async createMember(projectId: string, name: string, role: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/members`;
        const response = await this.http.post(path, JSON.stringify({
            name,
            role,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Requests members list. */
    public async getMembers(projectId: string): Promise<Member[]> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/members`;
        const response = await this.http.get(path);
        if (!response.ok) {
            await this.handleError(response);
        }
        const members = await response.json() || [new Member()];

        return members.map((member: Member) => new Member(member));
    };

    /** Updates member. */
    public async updateMember(member: Member): Promise<void> {
        const path = `${this.ROOT_PATH}/${member.projectId}/mint-site/members/${member.id}`;
        const response = await this.http.patch(path, JSON.stringify({
            name: member.name,
            role: member.role,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Updates member image. */
    public async updateMemberImage(projectId: string, memberId: string, image: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/members/${memberId}`;
        // TODO: string is valid JSON. So, it will be changed to  JSON.stringify(image)
        const response = await this.http.put(path, JSON.stringify({ image }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Deletes member. */
    public async deleteMember(projectId: string, memberId: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/members/${memberId}`;
        const response = await this.http.delete(path);
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Deletes member image. */
    public async deleteMemberImage(projectId: string, memberId: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/members/${memberId}/image`;
        const response = await this.http.delete(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Creates milestone. */
    public async createMilestone(projectId: string, milestoneFields: CreateMilestoneFields): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/milestones`;
        const response = await this.http.post(path, JSON.stringify(milestoneFields));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Requests milestones list. */
    public async getMilestones(projectId: string): Promise<Milestone[]> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/milestones`;
        const response = await this.http.get(path);
        if (!response.ok) {
            await this.handleError(response);
        }
        const milestones = await response.json() || [new Milestone()];

        return milestones.map((milestone: Milestone) => new Milestone(
            milestone.id,
            milestone.projectId,
            milestone.title,
            milestone.date,
            milestone.description,
        ));
    };

    /** Updates milestone. */
    public async updateMilestone(milestone: Milestone): Promise<void> {
        const path = `${this.ROOT_PATH}/${milestone.projectId}/mint-site/milestones/${milestone.id}`;
        const response = await this.http.patch(path, JSON.stringify({
            title: milestone.title,
            description: milestone.description,
            date: milestone.date,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Deletes milestone. */
    public async deleteMilestone(projectId: string, milestoneId: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/milestones/${milestoneId}`;
        const response = await this.http.delete(path);
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Creates social network. */
    public async createSocialNetwork(projectId: string, title: string, link: string) {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/social-networks`;
        const response = await this.http.post(path, JSON.stringify({
            title,
            link,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Requests social networks list. */
    public async getSocialNetworks(projectId: string): Promise<SocialNetwork[]> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/social-networks`;
        const response = await this.http.get(path);
        if (!response.ok) {
            await this.handleError(response);
        }
        const socialNetworks = await response.json() || [new SocialNetwork()];

        return socialNetworks.map((socialNetwork: SocialNetwork) => new SocialNetwork(
            socialNetwork.projectId,
            socialNetwork.id,
            socialNetwork.title,
            socialNetwork.link,
        ));
    };

    /** Updates social network. */
    public async updateSocialNetwork(socialNetwork: SocialNetwork): Promise<void> {
        const path = `${this.ROOT_PATH}/${socialNetwork.projectId}/mint-site/social-networks/${socialNetwork.id}`;
        const response = await this.http.patch(path, JSON.stringify(socialNetwork));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Deletes social network. */
    public async deleteSocialNetwork(projectId: string, socialNetworkId: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/social-networks/${socialNetworkId}`;
        const response = await this.http.delete(path);
        if (!response.ok) {
            await this.handleError(response);
        }
    };


    /** Gets project site 'content' images. */
    public async getContentImageNames(projectId: string): Promise<string[]> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/content`;
        const response = await this.http.get(path);
        if (!response.ok) {
            await this.handleError(response);
        }

        const imageNames = await response.json() || [];

        return imageNames.map((name: string) => name);
    };

    /** Updates project site 'content' image. */
    public async updateContentImage(projectId: string, image: string, name: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/content`;
        const response = await this.http.put(path, JSON.stringify({
            image,
            name: `${name}.png`,
        }));
        if (!response.ok) {
            await this.handleError(response);
        }
    };

    /** Updates project site 'content' image. */
    public async deleteContentImage(projectId: string, name: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/mint-site/content/${name}`;
        const response = await this.http.delete(path);
        if (!response.ok) {
            await this.handleError(response);
        }
    };
};
