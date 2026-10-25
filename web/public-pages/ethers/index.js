/** Functions of contracts/CollectionNFT.sol used by the mint site (ethers human-readable ABI). */
const COLLECTION_ABI = [
    'function mintWhitelist(bytes32[] _merkleProof, uint256 amount) payable',
    'function mint(uint256 numberOfTokens) payable',
    'function totalSupply() view returns (uint256)',
    'function presaleIsActive() view returns (bool)',
    'function saleIsActive() view returns (bool)',
    'function whitelistClaimed(address) view returns (bool)',
    'function MAX_SUPPLY() view returns (uint256)',
    'function PRICE() view returns (uint256)',
    'function MAX_PER_MINT() view returns (uint256)',
];

/** MetaMask error code: the requested chain has not been added to the wallet. */
const UNRECOGNIZED_CHAIN_CODE = 4902;
