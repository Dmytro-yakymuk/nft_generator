import { Mintsite, Project } from '@/collections';
import {
    SET_ABOUT_SECTION, SET_CONTENT_IMAGES, SET_GENERATION_COUNT,
    SET_HEADER, SET_MEMBER, SET_MEMBERS,
    SET_MILESTONE, SET_MILESTONES, SET_MINTSITE,
    SET_PROJECT, SET_PROJECTS, SET_SOCIAL_NETWORK,
    SET_SOCIAL_NETWORKS,
} from '@store/actions/collections';
import { ActionType } from '@/app/store';

const DEFAULT_VALUE = 0;

/** Exposes collections entity store state. Uses as default reducer state. */
class CollectionsState {
    /** At least, contains current project and all available projects. */
    constructor(
        public project: Project = new Project(),
        public projects: Project[] = [],
        public mintsite: Mintsite = new Mintsite(),
        public nftGenerationCount = DEFAULT_VALUE
    ) { };
};

export const collectionsReducer = (
    collectionsState: CollectionsState = new CollectionsState(),
    action: ActionType<any> = new ActionType('', ''),
) => {
    switch (action.type) {
    case SET_PROJECT:
        collectionsState.project = action.payload;
        break;
    case SET_PROJECTS:
        collectionsState.projects = action.payload;
        break;
    case SET_HEADER:
        collectionsState.mintsite.header = action.payload;
        break;
    case SET_ABOUT_SECTION:
        collectionsState.mintsite.aboutSection = action.payload;
        break;
    case SET_MINTSITE:
        collectionsState.mintsite = action.payload;
        break;
    case SET_MEMBER:
        collectionsState.mintsite.team[action.payload.index] = action.payload.member;
        collectionsState.mintsite.team = [...collectionsState.mintsite.team];
        break;
    case SET_MEMBERS:
        collectionsState.mintsite.team = action.payload;
        break;
    case SET_MILESTONE:
        collectionsState.mintsite.roadmap[action.payload.index] = action.payload.milestone;
        collectionsState.mintsite.roadmap = [...collectionsState.mintsite.roadmap];
        break;
    case SET_MILESTONES:
        collectionsState.mintsite.roadmap = action.payload;
        break;
    case SET_SOCIAL_NETWORK:
        collectionsState.mintsite.socialNetworks[action.payload.index] = action.payload.socialNetwork;
        collectionsState.mintsite.socialNetworks = [...collectionsState.mintsite.socialNetworks];
        break;
    case SET_SOCIAL_NETWORKS:
        collectionsState.mintsite.socialNetworks = action.payload;
        break;
    case SET_CONTENT_IMAGES:
        collectionsState.mintsite.contentImageNames = action.payload;
        break;
    case SET_GENERATION_COUNT:
        collectionsState.nftGenerationCount = action.payload;
        break;
    default:
        return collectionsState;
    };

    return { ...collectionsState };
};
