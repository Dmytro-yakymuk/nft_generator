/** Opens/closes metamask popup. */
function visibilityMetamaskPopup() {
    const popup = document.querySelector('.metamask-popup');
    const mintPopup = document.querySelector('.mint-popup');
    const popupWrapper = document.querySelector('.metamask-popup-wrapper');
    const mintPopupWrapper = document.querySelector('.mint-popup-wrapper');
    const targetClass = event.target.className;

    if (!popup || !popupWrapper) return;

    if (targetClass === 'mint') {
        popup.classList.remove('metamask-popup-hidden');
        popup.classList.add('metamask-popup-visible');
    }

    if (!popupWrapper.contains(event.target) && targetClass !== 'mint') {
        popup.classList.remove('metamask-popup-visible');
        popup.classList.add('metamask-popup-hidden');
    }

    // Clicks in the connect dialog open the mint dialog, so they must not close it,
    // whether the wallet answers before or after this handler runs.
    if (!mintPopupWrapper.contains(event.target) && !popupWrapper.contains(event.target) && targetClass !== 'mint-popup-btn') {
        mintPopup.classList.remove('mint-popup-visible');
        mintPopup.classList.add('mint-popup-hidden');
    }
};

/** Shows or hides element by toggling its "-visible" / "-hidden" classes. */
function setVisible(selector, visible) {
    const element = document.querySelector(`.${selector}`);

    if (!element) return;

    element.classList.toggle(`${selector}-visible`, visible);
    element.classList.toggle(`${selector}-hidden`, !visible);
};

/** Closes minted popup. */
function closeMintedPopup() {
    setVisible('minted-popup', false);
};

/** Sale terms of the contract (service.terms()), read when the page loads. */
let saleTerms = null;

/** Shows total price for the selected quantity. */
function updateMintPrice() {
    const mintCounter = document.querySelector('.mint-count');
    const mintCounterValue = document.querySelector('.mint-popup-info-value');

    if (!mintCounter || !mintCounterValue || !saleTerms) return;

    const total = saleTerms.price.mul(+mintCounter.textContent);
    mintCounterValue.textContent = `${ethers.utils.formatEther(total)} ETH`;
};

/** Decreases mint quantity number. */
function decreaseMintQuantity() {
    const mintCounter = document.querySelector('.mint-count');

    if (!mintCounter || +mintCounter.textContent <= 1) return;

    mintCounter.textContent--;
    updateMintPrice();
};

/** Increases mint quantity number up to the contract limit per transaction. */
function increaseMintQuantity() {
    const mintCounter = document.querySelector('.mint-count');

    const maxPerMint = saleTerms ? saleTerms.maxPerMint : window.MINT_SITE.maxPerMint;
    if (!mintCounter || +mintCounter.textContent >= maxPerMint) return;

    mintCounter.textContent++;
    updateMintPrice();
};

/** Shows minted tokens and the active sale phase from the contract. */
async function refreshSaleInfo() {
    const minted = document.querySelector('.mint-popup-minted');
    const phase = document.querySelector('.mint-popup-phase');

    try {
        const state = await service.saleState();
        if (minted) minted.textContent = `${state.minted}/${window.MINT_SITE.nftCount}`;
        if (phase) {
            phase.textContent = state.presale && state.sale ? 'Presale and public'
                : state.presale ? 'Presale (whitelist)'
                    : state.sale ? 'Public' : 'Not active';
        }
    } catch (error) {
        console.error('could not read sale state', error);
    }
};

/** Connects wallet */
async function connect () {
    if (!window.ethereum) {
        metamaskNotification(new MintError('Install MetaMask to mint'));

        return;
    }

    try {
        await window.ethereum.request({
            method: 'eth_requestAccounts',
        });
    } catch(error) {
        metamaskNotification(error);

        return;
    }

    setVisible('metamask-popup', false);
    setVisible('mint-popup', true);
    refreshSaleInfo();
}

/** Mints tokens: in the presale for whitelisted wallets, otherwise in the public sale. */
async function mint() {
    const mintCounter = document.querySelector('.mint-count');
    const mintedPopupInfo = document.querySelector('.minted-popup-info');
    const mintedPopupLink = document.querySelector('.minted-popup-link');

    if (!mintCounter || !+mintCounter.textContent) return;

    const count = +mintCounter.textContent;

    try {
        const receipt = await service.mint(count, () => {
            setVisible('mint-popup', false);
            setVisible('mint-loader', true);
        });

        setVisible('mint-loader', false);
        setVisible('minted-popup', true);
        mintedPopupInfo.textContent = `${count} ${count === 1 ? 'token' : 'tokens'} minted`;

        const link = service.transactionLink(receipt.transactionHash);
        mintedPopupLink.hidden = !link;
        mintedPopupLink.href = link;
    } catch(error) {
        setVisible('mint-loader', false);
        setVisible('mint-popup', false);
        metamaskNotification(error);
    }

    refreshSaleInfo();
}

document.addEventListener('DOMContentLoaded', () => {
    service.terms()
        .then((terms) => {
            saleTerms = terms;
            updateMintPrice();
        })
        .catch((error) => console.error('could not read sale terms', error));
    refreshSaleInfo();
});
