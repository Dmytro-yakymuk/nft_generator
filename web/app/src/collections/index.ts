/** Default number of photos. */
const DEFAULT_NUMBER_OF_NFTS: number = 0;

/** Defines project statuses.
 * Draft - project was created but not published.
 * None - project was not created.
 * Published - project was published.
*/
export enum ProjectStatus {
    'Draft' = 'Draft',
    'Published' = 'Published',
    'None' = 'None',
};

/** Defines project domain entity. */
export class Project {
    /** At least, contains project name. */
    constructor(
        public id: string = '1',
        public userId: string = '',
        public name: string = 'default',
        public contractAddress: string = '0x0000000000000000000000000000000000000000',
        public nftCount: number = DEFAULT_NUMBER_OF_NFTS,
        public isWebsitePublished: boolean = false,
        public areImagesGenerated: boolean = false,
        public websiteSlug: string = '',
        public imagesCid: string = '',
        public metadataCid: string = '',
        /** The pinning service that keeps the IPFS directories, empty while only the platform's node has them. */
        public ipfsPinnedIn: string = '',
    ) { };
};

/** Exposes website header section entity. */
export class Header {
    constructor(
        public projectId: string = '1',
        public title: string = 'Title',
        public description: string = 'About website',
        public logo: string = '',
        public image: string = '',
    ) {};
};

/** Exposes website about section entity. */
export class AboutSection {
    /** At least, contains projectId, title, and decsription. */
    constructor(
        public projectId: string = '1',
        public title: string = '',
        public description: string = '',
        public image: string = '',
    ) {};
};

/** Exposes website member entity. */
export class Member {
    public id: string = '';
    public projectId: string = 'default_project_id';
    public name: string = '';
    public role: string = '';
    public image: string = '';
    public isSaved: boolean = true;
    /** At least, contains id, projectId, name and role. */
    constructor(member?: Member) {
        Object.assign(this, member);
    };
};

/** Expose website milestone roadmap entity. */
export class Milestone {
    /** At least, contains projectId, id, title, date and description. */
    constructor(
        public id: string = '',
        public projectId: string = '',
        public title: string = '',
        public date: string = '',
        public description: string = '',
    ) { };
};

/** Exposes incoming fields for create new milestone. */
export class CreateMilestoneFields {
    /** At least, contains title, description and date. */
    constructor(
        public title: string = '',
        public description: string = '',
        public date: string = '',
    ) { };
};

/** Exposes social network entity. */
export class SocialNetwork {
    /** At least, contains projectId, id, title and link. */
    constructor(
        public projectId: string = '',
        public id: string = '',
        public title: string = '',
        public link: string = '',
    ) {};
};

/** Exposes project mintsite entity. */
export class Mintsite {
    /** At least, contains aboutSelection, team, roadmap and socialNetwork. */
    constructor(
        public header: Header = new Header(),
        public aboutSection: AboutSection = new AboutSection(),
        public team: Member[] = [],
        public roadmap: Milestone[] = [],
        public socialNetworks: SocialNetwork[] = [],
        public contentImageNames: string[] = [],
    ) {};
};
