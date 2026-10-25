/** Ethereum network the app works with; the server injects it into index.html. */
export interface Network {
    chainId: string;
    chainName: string;
    rpcUrl: string;
    blockExplorerUrl: string;
}

/** MetaMask error code: the requested chain has not been added to the wallet. */
const UNRECOGNIZED_CHAIN_CODE = 4902;

const DEFAULT_NETWORK: Network = {
    chainId: '0x7a69',
    chainName: 'Anvil',
    rpcUrl: 'http://localhost:8545',
    blockExplorerUrl: '',
};

/** Returns the configured network, or local Anvil if the page has no server config. */
export function getNetwork(): Network {
    return { ...DEFAULT_NETWORK, ...window.APP_CONFIG?.chain };
}

/** Switches MetaMask to the configured network, adding it to the wallet if needed. */
export async function switchNetwork(): Promise<void> {
    const network = getNetwork();
    const current: string = await window.ethereum.request({ method: 'eth_chainId' });
    if (current.toLowerCase() === network.chainId.toLowerCase()) {
        return;
    }

    try {
        await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: network.chainId }],
        });
    } catch (error: any) {
        if (error.code !== UNRECOGNIZED_CHAIN_CODE) {
            throw error;
        }
        await window.ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [{
                chainId: network.chainId,
                chainName: network.chainName,
                rpcUrls: [network.rpcUrl],
                nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
                blockExplorerUrls: network.blockExplorerUrl ? [network.blockExplorerUrl] : undefined,
            }],
        });
    }
}
