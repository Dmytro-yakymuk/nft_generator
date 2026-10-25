import { Link } from 'react-router-dom';

import { RoutesConfig } from '@/app/routes';
import { useMetamaskLogin } from '@/app/hooks/useMetamaskLogin';
import { useSession } from '@/app/hooks/useSession';

import './index.scss';

/** The last section of the home page (the Help link of the navbar): it leads to the studio, logging in with
 * MetaMask first if needed. The platform has no mail service, so there is no email form here. */
export const ContactUs = () => {
    const login = useMetamaskLogin();
    const isLoggedIn = useSession();

    return <section className="contact-us" id="contact-us">
        <span className="contact-us__title">Get Started</span>
        <span className="contact-us__subtitle">Let’s build something amazing.</span>
        <span className="contact-us__subtitle">Connect MetaMask, upload the layers and generate your collection.</span>
        {isLoggedIn
            ? <Link to={RoutesConfig.Collections.path} className="contact-us__start">
                <span className="contact-us__start__text">Open Studio</span>
            </Link>
            : <button className="contact-us__start" onClick={login}>
                <span className="contact-us__start__text">Connect MetaMask</span>
            </button>}
    </section>;
};
