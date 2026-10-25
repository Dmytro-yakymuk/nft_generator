import { useLocation } from 'react-router-dom';

import metamaskLogo from '@static/images/MetamaskModal/metamaskLogo.svg';

import { useMetamaskLogin } from '@/app/hooks/useMetamaskLogin';
import { useSession } from '@/app/hooks/useSession';

import './index.scss';

export const MetamaskModal = () => {
    const ROOT_PATH_LENGTH = 1;
    const location = useLocation();
    const login = useMetamaskLogin();
    const isLoggedIn = useSession();
    const isRootPath = location.pathname.length <= ROOT_PATH_LENGTH;

    return (
        <>
            {
                !isRootPath && !isLoggedIn ?
                    <div className="metamask-modal">
                        <div className="metamask-modal__wrapper">
                            <h3 className="metamask-modal__heading">
                                Connect  Metamask Wallet
                            </h3>
                            <p className="metamask-modal__description">
                                By connecting your wallet, you agree to our<br></br> Terms of Service and our Privacy Policy.
                            </p>
                            <img src={metamaskLogo} alt="metamaskLogo" />
                            <button className="metamask-modal__button" onClick={login}>
                                Connect Wallet
                            </button>
                        </div>
                    </div>
                    :
                    <></>
            }
        </>
    );
};
