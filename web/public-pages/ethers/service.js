/** Error with a message for the buyer. */
class MintError extends Error {}

class Service {
    client;
    site;
    termsRequest = null;

    constructor(ethersClient, site) {
        this.client = ethersClient;
        this.site = site;
    }

    /** Returns a read-only contract connected to the configured network. */
    readContract() {
        const provider = new ethers.providers.JsonRpcProvider(this.site.chain.rpcUrl);

        return new ethers.Contract(this.site.contractAddress, COLLECTION_ABI, provider);
    }

    /**
     * Returns the sale terms (collection size, price of one token in wei, limit per transaction) from the contract.
     * Contracts deployed before the terms became constructor parameters have no such functions: for them the terms
     * come with the page.
     */
    terms() {
        if (!this.termsRequest) {
            const contract = this.readContract();
            this.termsRequest = Promise.all([contract.MAX_SUPPLY(), contract.PRICE(), contract.MAX_PER_MINT()])
                .then(([maxSupply, price, maxPerMint]) => ({
                    maxSupply: maxSupply.toNumber(),
                    price,
                    maxPerMint: maxPerMint.toNumber(),
                }))
                .catch((error) => {
                    if (error.code !== 'CALL_EXCEPTION') {
                        this.termsRequest = null;
                        throw error;
                    }

                    return {
                        maxSupply: this.site.nftCount,
                        price: ethers.utils.parseEther(this.site.mintPrice),
                        maxPerMint: this.site.maxPerMint,
                    };
                });
        }

        return this.termsRequest;
    }

    /** Returns the number of minted tokens and the active sale phases. */
    async saleState() {
        const contract = this.readContract();
        const [minted, presale, sale] = await Promise.all([
            contract.totalSupply(),
            contract.presaleIsActive(),
            contract.saleIsActive(),
        ]);

        return { minted: minted.toNumber(), presale, sale };
    }

    /** Switches MetaMask to the configured network, adding it to the wallet if needed. */
    async switchNetwork() {
        const chain = this.site.chain;
        const current = await window.ethereum.request({ method: 'eth_chainId' });
        if (current.toLowerCase() === chain.chainId.toLowerCase()) {
            return;
        }

        try {
            await window.ethereum.request({
                method: 'wallet_switchEthereumChain',
                params: [{ chainId: chain.chainId }],
            });
        } catch (error) {
            if (error.code !== UNRECOGNIZED_CHAIN_CODE) {
                throw error;
            }
            await window.ethereum.request({
                method: 'wallet_addEthereumChain',
                params: [{
                    chainId: chain.chainId,
                    chainName: chain.chainName,
                    rpcUrls: [chain.rpcUrl],
                    nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
                    blockExplorerUrls: chain.blockExplorerUrl ? [chain.blockExplorerUrl] : undefined,
                }],
            });
        }
    }

    /**
     * Mints tokens for the connected wallet: in the presale with a Merkle proof if the wallet is
     * whitelisted and has not minted yet, otherwise in the public sale. Returns the transaction receipt.
     */
    async mint (mintCount, onSent) {
        await this.switchNetwork();

        // Network 'any' lets the provider follow the network MetaMask has just been switched to.
        const provider = new ethers.providers.Web3Provider(window.ethereum, 'any');
        const signer = provider.getSigner();
        const wallet = await signer.getAddress();
        const contract = new ethers.Contract(this.site.contractAddress, COLLECTION_ABI, signer);

        const [{ minted, presale, sale }, claimed, terms] = await Promise.all([
            this.saleState(),
            contract.whitelistClaimed(wallet),
            this.terms(),
        ]);
        if (minted + mintCount > terms.maxSupply) {
            throw new MintError(`Only ${Math.max(terms.maxSupply - minted, 0)} tokens left`);
        }
        if (mintCount > terms.maxPerMint) {
            throw new MintError(`At most ${terms.maxPerMint} tokens in one transaction`);
        }

        const value = terms.price.mul(mintCount);
        const proof = presale && !claimed ? await this.client.mercleProof(this.site.projectId, wallet) : null;

        let transaction;
        if (proof) {
            transaction = await contract.mintWhitelist(proof, mintCount, { value });
        } else if (sale) {
            transaction = await contract.mint(mintCount, { value });
        } else if (presale) {
            throw new MintError(claimed ? 'This wallet has already minted in the presale' : 'This wallet is not in the whitelist');
        } else {
            throw new MintError('Sale is not active');
        }

        onSent(transaction.hash);

        return await transaction.wait();
    }

    /** Returns a link to the transaction in the block explorer, or an empty string. */
    transactionLink(hash) {
        const explorer = this.site.chain.blockExplorerUrl;

        return explorer ? `${explorer.replace(/\/$/, '')}/tx/${hash}` : '';
    }
}
