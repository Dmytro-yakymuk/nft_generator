// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

contract CollectionNFT is ERC721, Ownable {
    // Sale terms the owner sets at deployment: the collection size, the price of one token, the limit per transaction.
    uint256 public immutable MAX_SUPPLY;
    uint256 public immutable PRICE;
    uint256 public immutable MAX_PER_MINT;

    string public PROVENANCE;
    bool public presaleIsActive = true;
    bool public saleIsActive = false;
    string private _baseURIextended;
    bytes32 public merkleRoot;
    uint public currentIndex = 0;
    mapping(address => bool) public whitelistClaimed;

    constructor(
        string memory baseURI,
        bytes32 merkleRoot_,
        string memory name_,
        uint256 maxSupply_,
        uint256 price_,
        uint256 maxPerMint_
    ) ERC721(name_, "NFT") {
        require(maxPerMint_ > 0 && maxPerMint_ <= maxSupply_, "Invalid limit per transaction");
        _baseURIextended = baseURI;
        merkleRoot = merkleRoot_;
        MAX_SUPPLY = maxSupply_;
        PRICE = price_;
        MAX_PER_MINT = maxPerMint_;
    }

    function setMerkleRoot(bytes32 merkleRoot_) external onlyOwner() {
        merkleRoot = merkleRoot_;
    }

    function setBaseURI(string memory baseURI_) external onlyOwner() {
        _baseURIextended = baseURI_;
    }

    function _baseURI() internal view virtual override returns (string memory) {
        return _baseURIextended;
    }

    function setProvenance(string memory provenance) public onlyOwner {
        PROVENANCE = provenance;
    }

    function setSaleState(bool newState) public onlyOwner {
        saleIsActive = newState;
    }

    function setPresaleState(bool newPresaleState) public onlyOwner {
        presaleIsActive = newPresaleState;
    }

    function totalSupply() view public returns (uint) {
        return currentIndex;
    }

    // Presale: a whitelisted wallet mints once with a Merkle proof of its address.
    function mintWhitelist(bytes32[] calldata _merkleProof, uint amount) public payable {
        require(presaleIsActive, "Presale must be active to mint");
        require(!whitelistClaimed[msg.sender], "Address has already claimed");
        require(msg.sender == tx.origin, "No transaction from smart contracts!");
        bytes32 leaf = keccak256(abi.encodePacked(msg.sender));
        require(MerkleProof.verify(_merkleProof, merkleRoot, leaf), "Invalid Proof");
        whitelistClaimed[msg.sender] = true;
        _mintTokens(amount);
    }

    // Public sale: any wallet mints within the limit per transaction.
    function mint(uint numberOfTokens) public payable {
        require(saleIsActive, "Sale must be active to mint");
        require(msg.sender == tx.origin, "No transaction from smart contracts!");
        _mintTokens(numberOfTokens);
    }

    // Checks the limits and the payment, then mints the next token ids to the caller.
    function _mintTokens(uint amount) private {
        require(amount > 0 && amount <= MAX_PER_MINT, "Too many tokens in one transaction");
        require(currentIndex + amount <= MAX_SUPPLY, "Purchase would exceed max supply");
        require(PRICE * amount <= msg.value, "Ether value sent is not correct");
        uint first = currentIndex;
        currentIndex = first + amount;
        for (uint i = 0; i < amount; i++) {
            _safeMint(msg.sender, first + i);
        }
    }

    // Withdraw the contract balance: the whole income goes to the owner of the collection.
    function withdraw() public onlyOwner {
        require(payable(msg.sender).send(address(this).balance));
    }
}
