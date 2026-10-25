/** Returns the contract revert reason from an ethers or MetaMask error, if any. */
function revertReason(error) {
    const nested = error.error || {};
    const message = error.reason || (nested.data && nested.data.message) || nested.message || (error.data && error.data.message) || '';
    const match = message.match(/execution reverted:?\s*(.*)/);

    return match ? (match[1] || 'Transaction reverted') : '';
}

/** Shows notification banner with custom text */
function metamaskNotification(error) {
    let errorMessage = '';
    const RPC_ERROR_CODE = -32002;
    const NOT_ENOUGH_ETH_CODE = 'INSUFFICIENT_FUNDS';
    const DENIED_TRANSACTION_CODE = 4001;
    const NETWORK_ERROR = 'NETWORK_ERROR';

    const metamaskNotificationElement = document.querySelector('.metamask-notification');

    if (error instanceof MintError) {
        errorMessage = error.message;
    } else {
        switch (error.code) {
        case RPC_ERROR_CODE:
            errorMessage = 'Please open metamask manually!';
            break;
        case DENIED_TRANSACTION_CODE:
            errorMessage = 'Transaction denied by user';
            break;
        case NOT_ENOUGH_ETH_CODE:
            errorMessage = 'Not enough balance';
            break;
        case NETWORK_ERROR:
            errorMessage = 'Network was changed, please refresh page';
            break;
        default:
            errorMessage = revertReason(error) || 'Something went wrong';
            break;
        }
    }
    console.error('[mint]', errorMessage, error);

    if (metamaskNotificationElement.style.display === 'none' || !metamaskNotificationElement.style.display) {
        metamaskNotificationElement.style.display = 'flex';
        metamaskNotificationElement.textContent = errorMessage;
        setTimeout(() => {
            metamaskNotificationElement.style.display = 'none';

        }, 4000)
    }
}
