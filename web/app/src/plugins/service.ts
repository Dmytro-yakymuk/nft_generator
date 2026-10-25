import { ethers } from 'ethers';

import { Service } from '@/ethers/service';

/** Creates Web3 provider dynamically to handle Metamask async injection */
export const getWeb3Provider = () => {
    if (!window.ethereum) {
        return null;
    }
    return new ethers.providers.Web3Provider(window.ethereum);
};

/** @deprecated Use getWeb3Provider() instead */
export const web3Provider = null;
export const ethersService = ethers;

/** Class for creating ethers service */
export class ServicePlugin {
    /** Creates ethers provider instance */
    public static create() {
        const provider = getWeb3Provider();
        return new Service(provider, ethers);
    }
}
