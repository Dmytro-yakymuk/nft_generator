import { UsersClient } from '@/api/users';
import { SignedMessage } from '@/ethers';
import { User, UserEditFields } from '@/accounts';

/**
 * Exposes all users related logic.
 */
export class UsersService {
    private readonly users: UsersClient;

    /** UsersService contains http implementation of users API. */
    public constructor(users: UsersClient) {
        this.users = users;
    };

    /** Returns user entity */
    public async getUser(): Promise<User> {
        return await this.users.getUser();
    }

    /** Updates user entity */
    public async updateuser(editFields: UserEditFields): Promise<void> {
        await this.users.updateUser(editFields);
    }

    /** Gets message from API for sign with metamask */
    public async getNonce(walletAddress: string): Promise<string> {
        const response = await this.users.getNonce(walletAddress);

        return response;
    }

    /** Gets message from API for sign with metamask */
    public async register(signature: string): Promise<void> {
        await this.users.register(signature);
    }

    /** Sends signed message, and logs-in */
    public async login(message: SignedMessage): Promise<void> {
        await this.users.login(message);
    }

    /** Logs out */
    public async logout(): Promise<void> {
        await this.users.logout();
    }
}
