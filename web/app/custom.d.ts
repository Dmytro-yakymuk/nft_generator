/** Ethereum interface */
interface ExternalProvider {
    isMetaMask?: boolean;
    isStatus?: boolean;
    host?: string;
    path?: string;
    chainId: string;
    sendAsync?: (request: { method: string; params?: any[] }, callback: (error: any, response: any) => void) => void;
    send?: (request: { method: string; params?: any[] }, callback: (error: any, response: any) => void) => void;
    request: (request: { method: string; params?: any[] }) => Promise<any>;
}

/** Runtime configuration injected by the server into index.html */
interface AppConfig {
    chain?: {
        chainId: string;
        chainName: string;
        rpcUrl: string;
        blockExplorerUrl: string;
    };
    /** HTTP gateway for files published to IPFS, e.g. http://localhost:8080 */
    ipfsGatewayUrl?: string;
    /** The service published collections can be pinned in, e.g. Filebase; empty if there is none. */
    ipfsPinning?: string;
    /** The domain whose subdomains serve published mint sites, e.g. localhost; empty if there is none. */
    mintSiteDomain?: string;
}

/** adds window.ethereum and window.APP_CONFIG field types */
declare global {
    interface Window {
        ethereum: ExternalProvider;
        APP_CONFIG?: AppConfig;
    }
}

export { };
