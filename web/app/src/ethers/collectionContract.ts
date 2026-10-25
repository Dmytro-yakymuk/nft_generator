import { ethers } from 'ethers';

import { getNetwork } from '@/ethers/network';

/** The functions of the collection contract (contracts/CollectionNFT.sol) the web app calls. */
const COLLECTION_ABI = [
    'function merkleRoot() view returns (bytes32)',
    'function setMerkleRoot(bytes32 merkleRoot_)',
    'function setBaseURI(string baseURI_)',
    'function totalSupply() view returns (uint256)',
    'function tokenURI(uint256 tokenId) view returns (string)',
];

const FIRST_TOKEN = 0;

/** Reads the whitelist Merkle root the collection contract holds through the RPC of the configured network,
 * so it does not depend on the network selected in MetaMask. */
export async function getContractMerkleRoot(contractAddress: string): Promise<string> {
    const provider = new ethers.providers.JsonRpcProvider(getNetwork().rpcUrl);
    const contract = new ethers.Contract(contractAddress, COLLECTION_ABI, provider);

    return await contract.merkleRoot();
}

/** Sets the whitelist Merkle root of the collection contract from the MetaMask wallet (only the contract owner can)
 * and waits until the transaction is mined. */
export async function setContractMerkleRoot(contractAddress: string, merkleRoot: string): Promise<void> {
    // Network 'any' lets the provider follow the network MetaMask has just been switched to.
    const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
    const contract = new ethers.Contract(contractAddress, COLLECTION_ABI, provider.getSigner());
    const transaction = await contract.setMerkleRoot(merkleRoot);

    await transaction.wait();
}

/** Returns the metadata address of the first token, or an empty string while no token is minted: the contract keeps
 * its base URI private, and tokenURI works only for an existing token. */
export async function getFirstTokenURI(contractAddress: string): Promise<string> {
    const provider = new ethers.providers.JsonRpcProvider(getNetwork().rpcUrl);
    const contract = new ethers.Contract(contractAddress, COLLECTION_ABI, provider);
    if ((await contract.totalSupply()).isZero()) {
        return '';
    }

    return await contract.tokenURI(FIRST_TOKEN);
}

/** Sets the base URI of the collection contract from the MetaMask wallet (only the contract owner can)
 * and waits until the transaction is mined. */
export async function setContractBaseURI(contractAddress: string, baseURI: string): Promise<void> {
    const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
    const contract = new ethers.Contract(contractAddress, COLLECTION_ABI, provider.getSigner());
    const transaction = await contract.setBaseURI(baseURI);

    await transaction.wait();
}
