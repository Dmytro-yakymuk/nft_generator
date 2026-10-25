import { ToastContainer, ToastPosition } from 'react-toastify';

import './index.scss';

/** Custom component for notifications. */
export const Notification: React.FC = () => {
    /** Describes notification position */
    const POSITION: ToastPosition = 'top-right';
    /** Closes notification after delay  */
    const AUTO_CLOSE_TIME: number = 5000;
    /** Describes notification queue */
    const IS_NEWEST_ON_TOP: boolean = false;
    /** Describes behaviour closes notification by onClick Mouse Event */
    const IS_CLOSED_ON_CLICK: boolean = false;
    /** RTL means right to left. Adds RTL support for pages */
    const IS_RIGHT_TO_LEFT_LAYOUT: boolean = false;

    return <ToastContainer
        position={POSITION}
        autoClose={AUTO_CLOSE_TIME}
        hideProgressBar
        newestOnTop={IS_NEWEST_ON_TOP}
        closeOnClick={IS_CLOSED_ON_CLICK}
        rtl={IS_RIGHT_TO_LEFT_LAYOUT}
        pauseOnFocusLoss
        pauseOnHover
    />;
};
