import { ethersService, getWeb3Provider } from '@/plugins/service';
import { ethers } from 'ethers';

/** Service for ethers methods */
export class Service {
    private provider: ethers.providers.Web3Provider | null;
    private readonly ethers;

    /** Applies ethereum provider for internal methods */
    public constructor(ethereumProvider: ethers.providers.Web3Provider | null, ethers: typeof ethersService) {
        this.provider = ethereumProvider;
        this.ethers = ethers;
    }

    /** Ensures provider is available, re-checking window.ethereum if needed */
    private ensureProvider(): ethers.providers.Web3Provider {
        if (!this.provider) {
            this.provider = getWeb3Provider();
        }
        if (!this.provider) {
            throw new Error('Metamask is not installed or not connected');
        }
        return this.provider;
    }

    /** Returns metamask wallet address */
    public async getWallet() {
        const signer = this.ensureProvider().getSigner();

        return await signer.getAddress();
    }

    /** Signs message and creates message raw signature */
    public async signMessage(message: string) {
        const signer = this.ensureProvider().getSigner();

        return await signer.signMessage(message);
    }

    /** Deploys contract with passed data and waits until the transaction is mined */
    public async deployContract(data: string): Promise<ethers.providers.TransactionReceipt> {
        // Network 'any' lets the provider follow the network MetaMask has just been switched to.
        const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
        const transaction = await provider.getSigner().sendTransaction({
            data: `0x${data}`,
        });

        return await transaction.wait();
    }
}
