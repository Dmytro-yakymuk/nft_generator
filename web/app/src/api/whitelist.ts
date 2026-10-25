import { APIClient } from '@/api';
import { WhitelistPage, WhitelistedWallet } from '@/whitelist';

/**
 * WhitelistClient is a http implementation of whitelist API.
 * Exposes all whitelist-related functionality.
 */
export class WhitelistClient extends APIClient {
    private readonly ROOT_PATH: string = '/api/v0/whitelist';

    /** Requests whitelist wallets. */
    public async getWhitelist(projectId: string, currentPage: number, name: string): Promise<WhitelistPage> {
        /** Wallets count on page. */
        const limit: number = 10;
        const path = `${this.ROOT_PATH}/${projectId}?page=${currentPage}&limit=${limit}&name=${name}`;

        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const whitelistPage = await response.json();

        return new WhitelistPage(
            whitelistPage.items && whitelistPage.items.map(
                (whiteListedWallet: WhitelistedWallet) => new WhitelistedWallet(whiteListedWallet.name, whiteListedWallet.walletAddress)),
            whitelistPage.page
        );
    };

    /** Addes wallet address to whitelist. */
    public async addWalletToWhitelist(projectId: string, name: string, walletAddress: string): Promise<void> {
        const path = `${this.ROOT_PATH}`;
        const response = await this.http.post(path, JSON.stringify({ projectId, name, walletAddress }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Deletes wallet address from whitelist. */
    public async deleteWalletFromWhitelist(projectId: string, walletAddress: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/${walletAddress}`;
        const response = await this.http.delete(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Returns the Merkle root of the whitelist that the deployed contract must hold (zero for an empty whitelist). */
    public async getMerkleRoot(projectId: string): Promise<string> {
        const path = `${this.ROOT_PATH}/${projectId}/merkle-root`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const { merkleRoot } = await response.json();

        return merkleRoot;
    }

    /** Uploads whitelist excel file. */
    public async uploadWhitelistFile(projectId: string, file: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/excel-file`;
        const response = await this.http.post(path, JSON.stringify({ file }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }
}
