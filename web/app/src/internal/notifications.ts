import { toast } from 'react-toastify';

/** Code which indicates that 'eth_requestAccounts' already processing */
const RPC_ERROR_CODE = -32002;
const ALREADY_MINTED_ERROR_CODE = -32603;
const DENIED_TRANSACTION_CODE = 4001;
const NOT_ENOUGH_ETH_CODE = 'INSUFFICIENT_FUNDS';
const NETWORK_ERROR = 'NETWORK_ERROR';

/** Handles metamask errors */
export function metamaskError(error: any) {
    let errorMessage = '';
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
        errorMessage = 'Something went wrong';
        break;
    }
    toast.error(errorMessage, {
        position: toast.POSITION.TOP_RIGHT,
        theme: 'colored',
    });
}

/** Toast types */
type ToastType = 'error' | 'info' | 'success' | 'warning';
/** Toast themes type */
type ToastTheme = 'colored' | 'light' | 'dark';

/** Function to make toast notification syntax shotrer */
export function notify(message: string, type: ToastType = 'error', theme: ToastTheme = 'colored') {
    toast[type](message, {
        position: toast.POSITION.TOP_RIGHT,
        theme: theme,
    });
}
