import { Dispatch } from 'redux';

import { ProjectsClient } from '@/api/projects';
import { AboutSection, CreateMilestoneFields, Header, Member, Milestone, Mintsite, Project, SocialNetwork } from '@/collections';
import { ProjectsService } from '@/collections/service';

/** SET_PROJECT is action type for direct manipulation setProject action of projects store. */
export const SET_PROJECT: string = 'SET_PROJECT';
/** SET_PROJECTS is action type for direct manipulation setProjects action of projects store. */
export const SET_PROJECTS: string = 'SET_PROJECTS';
/** SET_HEADER is action type for direct manipulation setHeader action of projects store. */
export const SET_HEADER = 'SET_HEADER';
/** SET_ABOUT_SECTION is action type for direct manipulation setAboutSection action of projects store. */
export const SET_ABOUT_SECTION: string = 'SET_ABOUT_SECTION';
/** SET_MINTSITE is action type for direct manipulation setMintsite action of projects store. */
export const SET_MINTSITE: string = 'SET_MINTSITE';
/** SET_MEMBER is action type for direct manipulation setMember action of projects store. */
export const SET_MEMBER: string = 'SET_MEMBER';
/** SET_MEMBERS is action type for direct manipulation setMembers action of projects store. */
export const SET_MEMBERS: string = 'SET_MEMBERS';
/** SET_MILESTONE is action type for direct manipulation setMilestone action of projects store. */
export const SET_MILESTONE: string = 'SET_MILESTONE';
/** SET_MILESTONES is action type for direct manipulation setMilestones action of projects store. */
export const SET_MILESTONES: string = 'SET_MILESTONES';
/** SET_SOCIAL_NETWORK is action type for direct manipulation setSocialNetworks action of projects store. */
export const SET_SOCIAL_NETWORK: string = 'SET_SOCIAL_NETWORK';
/** SET_SOCIAL_NETWORKS is action type for direct manipulation setSocialNetworks action of projects store. */
export const SET_SOCIAL_NETWORKS: string = 'SET_SOCIAL_NETWORKS';
/** SET_CONTENT_IMAGES is action type for direct manipulation setContentImages action of projects store. */
export const SET_CONTENT_IMAGES = 'SET_CONTENT_IMAGES';
/** SET_GENERATION_COUNT is action type for direct manipulation setGenerationProgres action of projects store. */
export const SET_GENERATION_COUNT = 'SET_GENERATION_COUNT';

/** An action setProject contains type and payload data for sets project. */
const setProject = (project: Project) => ({
    type: SET_PROJECT,
    payload: project,
});

/** An action setProjects contains type and payload data for sets projects. */
const setProjects = (projects: Project[]) => ({
    type: SET_PROJECTS,
    payload: projects,
});

/** An action setAboutSection contains type and payload data for sets about section. */
export const setAboutSection = (aboutSection: AboutSection) => ({
    type: SET_ABOUT_SECTION,
    payload: aboutSection,
});

/** An action setHeader contains type and payload data for sets header section. */
export const setHeader = (header: Header) => ({
    type: SET_HEADER,
    payload: header,
});

/** An action setProjectMintsite contains type and payload data for sets mintsite. */
const setMintsite = (mintsite: Mintsite) => ({
    type: SET_MINTSITE,
    payload: mintsite,
});

/** An action setMember contains type and payload data for sets project member. */
export const setMember = (index: number, member: Member) => ({
    type: SET_MEMBER,
    payload: {
        index,
        member,
    },
});

/** An action setMembers contains type and payload data for sets project members. */
export const setMembers = (members: Member[]) => ({
    type: SET_MEMBERS,
    payload: members,
});

/** An action setMilestone contains type and payload data for sets roadmap milestone. */
export const setMilestone = (index: number, milestone: Milestone) => ({
    type: SET_MILESTONE,
    payload: {
        index,
        milestone,
    },
});

/** An action setMilestones contains type and payload data for sets project milestones. */
export const setMilestones = (milestones: Milestone[]) => ({
    type: SET_MILESTONES,
    payload: milestones,
});

/** An action socialNetwork contains type and payload data for sets project social network. */
export const setSocialNetwork = (index: number, socialNetwork: SocialNetwork) => ({
    type: SET_SOCIAL_NETWORK,
    payload: {
        index,
        socialNetwork,
    },
});

/** An action setSocialNetworks contains type and payload data for sets project social networks. */
export const setSocialNetworks = (socialNetworks: SocialNetwork[]) => ({
    type: SET_SOCIAL_NETWORKS,
    payload: socialNetworks,
});

/** An action setContentImages contains type and payload data for sets project content images. */
export const setContentImages = (images: string[]) => ({
    type: SET_CONTENT_IMAGES,
    payload: images,
});
/** An action setGenerationProgres contains type and payload data for sets project content images. */
export const setGenerationProgress = (count: number) => ({
    type: SET_GENERATION_COUNT,
    payload: count,
});

const projectsClient = new ProjectsClient();
const projectsService = new ProjectsService(projectsClient);

/** Thunk middleware that creates new project and sets updated projects list to store. */
export const createProject = (name: string, nftCount: number) => async function(dispatch: Dispatch) {
    await projectsService.create(name, nftCount);
    const projects = await projectsService.list();
    dispatch(setProjects(projects));
};

/** Thunk middleware that deletes project and sets updated projects list to store. */
export const deleteProject = (id: string) => async function(dispatch: Dispatch) {
    await projectsService.delete(id);
    const projects = await projectsService.list();
    dispatch(setProjects(projects));
};

/** Thunk middleware that updates project and sets updated projects list to store. */
export const updateProject = (id: string, project: Project) => async function(dispatch: Dispatch) {
    await projectsService.updateProject(id, project);
    const projects = await projectsService.list();
    dispatch(setProjects(projects));
};

/** Thunk middleware that requests project and set project to store. */
export const getProject = (id: string) => async function(dispatch: Dispatch) {
    const project = await projectsService.getProject(id);
    dispatch(setProject(project));
};

/** Saves the deployed contract address and the number of the block it was deployed in. */
export const updateContractAddress = (projectId: string, address: string, block: number) => async(dispatch: Dispatch) => {
    await projectsService.updateContractAddress(projectId, address, block);
};

/** Publishes project mint site at /sites/{slug} and refreshes the project. */
export const publishWebsite = (projectId: string, slug: string) => async function(dispatch: Dispatch) {
    await projectsService.publishWebsite(projectId, slug);
    const project = await projectsService.getProject(projectId);
    dispatch(setProject(project));
};

/** Thunk middleware that requests all projects and sets projects to store. */
export const getProjects = () => async function(dispatch: Dispatch) {
    const projects = await projectsService.list();
    dispatch(setProjects(projects));
};

/** Sends request to the api for NFT generating */
export const generateNfts = (id: string) => async function(dispatch: Dispatch) {
    await projectsService.generateNfts(id);
};

/** Deletes generated NFT's. */
export const deleteGeneratedNFTs = async(id: string): Promise<void> => await projectsService.deleteGeneratedNFTs(id);

/** Aborts NFT's generation. */
export const abortGenerateNTTs = async(id: string):Promise<void> => await projectsService.abortGenerateNFTs(id);

/** Returns progress of current nft generation */
export const getGenerationProgress = (id: string) => async function(dispatch: Dispatch) {
    const count = await projectsService.getGenerationProgress(id);
    dispatch(setGenerationProgress(count));
};

/** Generates metadata for the generated collection. */
export const generateCollectionMetadata = async(id: string):Promise<void> => await projectsService.generateCollectionMetadata(id);

/** Thunk middleware that publishes the generated collection to IPFS and sets the updated project to store. */
export const publishToIPFS = (id: string) => async function(dispatch: Dispatch) {
    const project = await projectsService.publishToIPFS(id);
    dispatch(setProject(project));
};

/** Thunk middleware that pins the published collection in the pinning service and sets the updated project to store. */
export const pinIPFS = (id: string) => async function(dispatch: Dispatch) {
    const project = await projectsService.pinIPFS(id);
    dispatch(setProject(project));
};

/** Thunk middleware that requests project mintsite and sets project mintsite to store. */
export const getProjectMintsite = (projectId: string) => async function(dispatch: Dispatch) {
    const mintsite = await projectsService.getProjectMintSite(projectId);
    dispatch(setMintsite(mintsite));
};

/** Deleletes project mintsite. */
export const deleteMintsite = async(projectId: string): Promise<void> => await projectsService.deleteMintSite(projectId);

/** Thunk middleware that creates new 'header' section and sets updated section list to store. */
export const createHeader = (projectId: string, title: string, decsription: string) => async function(dispatch: Dispatch) {
    await projectsService.createHeader(projectId, title, decsription);
    const header = await projectsService.getHeader(projectId);
    dispatch(setHeader(header));
};

/** Thunk middleware that requests project mintsite and sets project mintsite to store. */
export const getHeader = (projectId: string) => async function(dispatch: Dispatch) {
    const header = await projectsService.getHeader(projectId);
    dispatch(setHeader(header));
};

/** Thunk middleware that updates 'header' text fields and sets updated section list to store. */
export const updateHeaderTextFields = (header: Header) => async function(dispatch: Dispatch) {
    await projectsService.updateHeaderTextFields(header);
    const updatedHeader = await projectsService.getHeader(header.projectId);
    dispatch(setHeader(updatedHeader));
};

/** Thunk middleware that updates 'header' logo and sets updated section list to store. */
export const updateHeaderLogo = (header: Header) => async function(dispatch: Dispatch) {
    await projectsService.updateHeaderLogo(header);
    const updatedHeader = await projectsService.getHeader(header.projectId);
    dispatch(setHeader(updatedHeader));
};

/** Thunk middleware that updates 'header' image and sets updated section list to store. */
export const updateHeaderImage = (header: Header) => async function(dispatch: Dispatch) {
    await projectsService.updateHeaderImage(header);
    const updatedHeader = await projectsService.getHeader(header.projectId);
    dispatch(setHeader(updatedHeader));
};

/** Thunk middleware that requests project mintsite and sets project mintsite to store. */
export const getAboutSection = (projectId: string) => async function(dispatch: Dispatch) {
    const aboutSection = await projectsService.getAboutSection(projectId);
    dispatch(setAboutSection(aboutSection));
};

/** Thunk middleware that creates new 'about' section and sets updated section list to store. */
export const createAboutSection = (projectId: string, title: string, decsription: string) => async function(dispatch: Dispatch) {
    await projectsService.createAboutSection(projectId, title, decsription);
    const aboutSection = await projectsService.getAboutSection(projectId);
    dispatch(setAboutSection(aboutSection));
};

/** Thunk middleware that updates 'about' section and sets updated section to store. */
export const updateAboutSection = (aboutSection: AboutSection) => async function(dispatch: Dispatch) {
    await projectsService.updateAboutSection(aboutSection);
    await projectsService.updateAboutSectionImage(aboutSection.projectId, aboutSection.image);
    const updatedAboutSection = await projectsService.getAboutSection(aboutSection.projectId);
    dispatch(setAboutSection(updatedAboutSection));
};

/** Thunk middleware that creates new member and sets updates members list to store. */
export const createMember = (projectId: string, name: string, role: string) => async function(dispatch: Dispatch) {
    await projectsService.createMember(projectId, name, role);
};

/** Thunk middleware that requests members list and sets members list to store. */
export const getMembers = (projectId: string) => async function(dispatch: Dispatch) {
    const members = await projectsService.getMembers(projectId);
    dispatch(setMembers(members));
};

/** Thunk middleware that updates member. */
export const updateMember = (member: Member) => async function(dispatch: Dispatch) {
    await projectsService.updateMember(member);
};

/** Updates member image. */
export const updateMemberImage = (projectId: string, memberId: string, image: string) => async function(dispatch: Dispatch) {
    await projectsService.updateMemberImage(projectId, memberId, image);
};

/** Deletes member image. */
export const deleteMemberImage = (projectId: string, memberId: string) => async function(dispatch: Dispatch) {
    await projectsService.deleteMemberImage(projectId, memberId);
};

/** Thunk middleware that deletes member and updates members list to store. */
export const deleteMember = (projectId: string, memberId: string) => async function(dispatch: Dispatch) {
    await projectsService.deleteMember(projectId, memberId);
};

/** Thunk middleware that creates new milestone and sets updated milestones list to store. */
export const createMilestone = (projectId: string, milestoneFields: CreateMilestoneFields) => async function(dispatch: Dispatch) {
    await projectsService.createMilestone(projectId, milestoneFields);
};

/** Thunk middleware that requests milestones list and sets members list to store. */
export const getMilestones = (projectId: string) => async function(dispatch: Dispatch) {
    const milestones = await projectsService.getMilestones(projectId);
    dispatch(setMilestones(milestones));
};

/** Thunk middleware that updates milestone and sets updated milestones list to store. */
export const updateMilestone = (milestone: Milestone) => async function(dispatch: Dispatch) {
    await projectsService.updateMilestone(milestone);
};

/** Thunk middleware that deletes milestone and updates milestones list to store. */
export const deleteMilestone = (projectId: string, milestoneId: string) => async function(dispatch: Dispatch) {
    await projectsService.deleteMilestone(projectId, milestoneId);
};

/** Thunk middleware that creates new social network and sets updated social networks list to store. */
export const createSocialNetwork = (projectId: string, title: string, link: string) => async function(dispatch: Dispatch) {
    await projectsService.createSocialNetwork(projectId, title, link);
};

/** Thunk middleware that requests social networks list and sets social networks list to store. */
export const getSocialNetworks = (projectId: string) => async function(dispatch: Dispatch) {
    const socialNetworks = await projectsService.getSocialNetworks(projectId);
    dispatch(setSocialNetworks(socialNetworks));
};

/** Thunk middleware that updates social network and sets updated social networks list to store. */
export const updateSocialNetwork = (socialNetwork: SocialNetwork) => async function(dispatch: Dispatch) {
    await projectsService.updateSocialNetwork(socialNetwork);
};

/** Thunk middleware that deletes social network and updates social networks list to store. */
export const deleteSocialNetwork = (projectId: string, socialNetworkId: string) => async function(dispatch: Dispatch) {
    await projectsService.deleteSocialNetwork(projectId, socialNetworkId);
};

/** Thunk middleware that requests content image`s names and sets names to store. */
export const getContentImageNames = (projectId: string) => async function(dispatch: Dispatch) {
    const imageNames = await projectsService.getContentImageNames(projectId);
    dispatch(setContentImages(imageNames));
};

/** Thunk middleware that updates content image and sets updated images to store. */
export const updateContentImage = (projectId: string, image: string, name: string) => async function(dispatch: Dispatch) {
    await projectsService.updateContentImage(projectId, image, name);
    const imageNames = await projectsService.getContentImageNames(projectId);
    dispatch(setContentImages(imageNames));
};

/** Thunk middleware that deletes content image and sets updated images to store. */
export const deleteContentImage = (projectId: string, name: string) => async function(dispatch: Dispatch) {
    await projectsService.deleteContentImage(projectId, name);
    const imageNames = await projectsService.getContentImageNames(projectId);
    dispatch(setContentImages(imageNames));
};
