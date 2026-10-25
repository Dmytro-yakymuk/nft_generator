import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import { NavbarSettings } from '@views/Home';

import { RoutesConfig } from '@/app/routes';
import { session } from '@/accounts/session';
import { useMetamaskLogin } from '@/app/hooks/useMetamaskLogin';
import { useSession } from '@/app/hooks/useSession';

import './index.scss';

/** Class that describes the entity of the navbar link. */
class NavbarLink {
    public constructor(
        public name: string,
        public link: string,
    ) { }
}

export const Navbar: React.FC<{ navbarSettings: NavbarSettings }> = ({ navbarSettings }) => {
    const FIRST_ELEMENT_INDEX = 0;
    const FIFTH_ELEMENT_INDEX = 4;
    const FOURTH_ELEMENT_FROM_END_INDEX = -4;
    const [walletLabel, setWalletLabel] = useState('');
    const login = useMetamaskLogin();
    const isLoggedIn = useSession();

    /** Navbar links: the studio and the sections of the home page. */
    const navbarLinks: NavbarLink[] = [
        new NavbarLink('Studio', RoutesConfig.Collections.path),
        new NavbarLink('Creators', '/#how-it-works'),
        new NavbarLink('Brands', '/#featured-drop'),
        new NavbarLink('Help', '/#contact-us'),
        new NavbarLink('About', '/#features'),
    ];

    useEffect(() => {
        (async function() {
            // the address comes from the wallet itself: the navbar of every page needs no ethers.js.
            try {
                const [address]: string[] = window.ethereum ? await window.ethereum.request({ method: 'eth_accounts' }) : [];
                if (!address) {
                    throw new Error('no connected account');
                }
                setWalletLabel(`${address.slice(FIRST_ELEMENT_INDEX, FIFTH_ELEMENT_INDEX)}...${address.slice(FOURTH_ELEMENT_FROM_END_INDEX)}`);
            } catch (error: any) {
                session.end();
            }
        })();
    }, []);

    return (
        <nav className={`navbar ${navbarSettings.borderColor} ${navbarSettings.backgroundColor}`}>
            <div className="wrapper">
                <a href="/" className="navbar__logo">
                    <img src={navbarSettings.logo} alt="Generator logo" />
                </a>
                <ul className="navbar__items">
                    {
                        navbarLinks.map(
                            (item: NavbarLink, index: number) =>
                                <li className={`navbar__items__item-${navbarSettings.textColor}`} key={index}>
                                    <Link to={item.link}>{item.name}</Link>
                                </li>
                        )
                    }
                </ul>
                {
                    isLoggedIn
                        ?
                        <Link
                            to={RoutesConfig.AccountSettings.path}
                            className={`navbar__connect-btn ${navbarSettings.borderColor}`}
                        >
                            <span className={`navbar__connect-btn__${navbarSettings.textColor}`}>{walletLabel}</span>
                        </Link>
                        :
                        <button
                            className={`navbar__connect-btn ${navbarSettings.borderColor}`}
                            onClick={login}
                        >
                            <span className={`navbar__connect-btn__${navbarSettings.textColor}`}>Connect</span>
                        </button>
                }
            </div>
        </nav>
    );
};
