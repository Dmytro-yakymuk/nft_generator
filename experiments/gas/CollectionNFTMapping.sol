// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// Experiment E5 (section 3 of the thesis): CollectionNFT with the whitelist kept as a mapping of
// addresses instead of a Merkle root. Everything else, including minting, is the same as in CollectionNFT.
contract CollectionNFTMapping is ERC721, Ownable {
    uint256 public immutable MAX_SUPPLY;
    uint256 public immutable PRICE;
    uint256 public immutable MAX_PER_MINT;

    string public PROVENANCE;
    bool public presaleIsActive = true;
    bool public saleIsActive = false;
    string private _baseURIextended;
    mapping(address => bool) public whitelist;
    uint public currentIndex = 0;
    mapping(address => bool) public whitelistClaimed;

    constructor(
        string memory baseURI,
        string memory name_,
        uint256 maxSupply_,
        uint256 price_,
        uint256 maxPerMint_
    ) ERC721(name_, "NFT") {
        require(maxPerMint_ > 0 && maxPerMint_ <= maxSupply_, "Invalid limit per transaction");
        _baseURIextended = baseURI;
        MAX_SUPPLY = maxSupply_;
        PRICE = price_;
        MAX_PER_MINT = maxPerMint_;
    }

    function addToWhitelist(address[] calldata addresses) external onlyOwner() {
        for (uint i = 0; i < addresses.length; i++) {
            whitelist[addresses[i]] = true;
        }
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

    function mintWhitelist(uint amount) public payable {
        require(presaleIsActive, "Presale must be active to mint");
        require(!whitelistClaimed[msg.sender], "Address has already claimed");
        require(msg.sender == tx.origin, "No transaction from smart contracts!");
        require(whitelist[msg.sender], "Not in the whitelist");
        whitelistClaimed[msg.sender] = true;
        _mintTokens(amount);
    }

    function mint(uint numberOfTokens) public payable {
        require(saleIsActive, "Sale must be active to mint");
        require(msg.sender == tx.origin, "No transaction from smart contracts!");
        _mintTokens(numberOfTokens);
    }

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

    function withdraw() public onlyOwner {
        require(payable(msg.sender).send(address(this).balance));
    }
}
