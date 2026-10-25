import { APIClient } from '@/api';
import { SignedMessage } from '@/ethers';
import { User, UserEditFields } from '@/accounts';

/**
 * UsersClient is a http implementation of users API.
 * Exposes all users-related functionality.
 */
export class UsersClient extends APIClient {
    private readonly ROOT_PATH: string = '/api/v0';

    /** Returns user entity */
    public async getUser(): Promise<User> {
        const path = `${this.ROOT_PATH}/accounts`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        };

        const user: User = await response.json();

        return Object.assign(new User(), user);
    }

    /** Updates user entity */
    public async updateUser(editFields: UserEditFields): Promise<void> {
        const path = `${this.ROOT_PATH}/accounts`;
        const response = await this.http.patch(path, JSON.stringify(editFields));

        if (!response.ok) {
            await this.handleError(response);
        };
    }

    /** Gets message from API for sign with metamask */
    public async getNonce(walletAddress: string): Promise<string> {
        const response = await this.http.get(`${this.ROOT_PATH}/auth/metamask/nonce?address=${walletAddress}`);

        if (!response.ok) {
            await this.handleError(response);
        };

        return await response.json();
    }

    /** Sends signed message and registers user */
    public async register(signature: string): Promise<void> {
        const response = await this.http.post(`${this.ROOT_PATH}/auth/metamask/register`, JSON.stringify(signature));

        if (!response.ok) {
            await this.handleError(response);
        };
    }

    /** Sends signed message, and logs-in */
    public async login(message: SignedMessage): Promise<void> {
        const response = await this.http.post(`${this.ROOT_PATH}/auth/metamask/login`, JSON.stringify(message));

        if (!response.ok) {
            await this.handleError(response);
        };
    }

    /** Logs out: the server removes the auth cookie */
    public async logout(): Promise<void> {
        const response = await this.http.post(`${this.ROOT_PATH}/auth/logout`);

        if (!response.ok) {
            await this.handleError(response);
        };
    }
}
