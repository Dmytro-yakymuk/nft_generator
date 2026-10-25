/** Whitelisted wallet address entity. */
export class WhitelistedWallet {
    public constructor(
        public name: string = '',
        public walletAddress: string = '',
    ) {}
};


/** Defines whitelist page pagination. */
export class WhitelistPage {
    constructor(
        public wallets: WhitelistedWallet[],
        public page: {
            offset: number;
            limit: number;
            currentPage: number;
            pageCount: number;
            totalCount: number;
        }
    ) {};
};
