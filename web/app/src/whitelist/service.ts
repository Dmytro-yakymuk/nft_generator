import { WhitelistClient } from '@/api/whitelist';
import { WhitelistPage } from '@/whitelist';

/** Default selected page.
 * Uses as defalt parameter if selected page is not given.
 */
const DEFAULT_SELECTED_PAGE: number = 1;
const DEFAULT_WALLET_NAME: string = '';
export class WhitelistService {
    private readonly whitelist: WhitelistClient;

    constructor(whitelist: WhitelistClient) {
        this.whitelist = whitelist;
    }

    /** Requests whitelist. */
    public async getWhitelist(projectId: string, currentPage: number = DEFAULT_SELECTED_PAGE, name: string = DEFAULT_WALLET_NAME): Promise<WhitelistPage> {
        return await this.whitelist.getWhitelist(projectId, currentPage, name);
    };

    /** Addes wallet address to whitelist. */
    public async addWalletToWhitelist(projectId: string, name: string, walletAddress: string): Promise<void> {
        await this.whitelist.addWalletToWhitelist(projectId, name, walletAddress);
    }

    /** Deletes wallet address from whitelist. */
    public async deleteWalletFromWhitelist(projectId: string, walletAddress: string): Promise<void> {
        await this.whitelist.deleteWalletFromWhitelist(projectId, walletAddress);
    }

    /** Returns the Merkle root of the whitelist that the deployed contract must hold. */
    public async getMerkleRoot(projectId: string): Promise<string> {
        return await this.whitelist.getMerkleRoot(projectId);
    }

    /** Uploads whitelist excel file. */
    public async uploadWhitelistFile(projectId: string, file: string): Promise<void> {
        await this.whitelist.uploadWhitelistFile(projectId, file);
    }
}
