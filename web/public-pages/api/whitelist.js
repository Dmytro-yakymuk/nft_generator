/** Mint site API client; the routes do not require authentication. */
class EthersClient extends APIClient {
    ROOT_PATH = '/api/v0/public';

    /** Returns Merkle proof of the wallet address, or null if the address is not in the whitelist. */
    async mercleProof (projectId, walletAddress) {
        const response = await this.http.get(`${this.ROOT_PATH}/whitelist/${projectId}/${walletAddress}/merkle-proof`);

        if (response.status === NOT_FOUND_ERROR) {
            return null;
        }
        if (!response.ok) {
            await this.handleError(response);
        }

        return await response.json();
    }
}
