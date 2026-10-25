import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

import { Navbar } from '@components/common/Navbar';
import { MintedBy } from '@components/home/MintedBy';
import { FeaturedDrop } from '@components/home/FeaturedDrop';
import { ContactUs } from '@components/home/ContactUs';
import { HowItWorks } from '@components/home/HowItWorks';

import pinkLogo from '@static/images/logo.png';
import whiteLogo from '@static/images/secondLogo.png';
import blackLogo from '@static/images/thirdLogo.png';
import fourthLogo from '@static/images/howItWorks/logo_white.png';

/** Class that describes navbar styles. */
export class NavbarSettings {
    public constructor(
        public borderColor: string,
        public textColor: string,
        public backgroundColor: string,
        public logo: string,
    ) {}
}

/** Maximum height to first home section. */
const HEIGHT_FROM_TOP_TO_FIRST_SECTION: number = 1692;
const HEIGHT_FROM_TOP_TO_THIRD_SECTION: number = 124;
const HEIGHT_FROM_TOP_TO_FOURTH_SECTION: number = -960;

const Home = () => {
    const [navbarSettings, setNavbarSettings] = useState<NavbarSettings>(new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo));

    const changeNavbarStyles = () => {
        const firstHomeSection = document.querySelector('.minted-by');
        const secondHomeSection = document.querySelector('.featured-drop');
        const thirdHomeSection = document.querySelector('.contact-us');
        const fourthHomeSection = document.querySelector('.works');
        const navbar = document.querySelector('.navbar');

        if (!firstHomeSection || !navbar || !secondHomeSection || !thirdHomeSection || !fourthHomeSection) {
            return;
        }

        /** Current height from top page to home sections. */
        const heightToFirstSection = thirdHomeSection.getBoundingClientRect().top;
        const heightToSecondSection = secondHomeSection.getBoundingClientRect().top;

        if (secondHomeSection instanceof HTMLElement && heightToSecondSection >= HEIGHT_FROM_TOP_TO_FIRST_SECTION) {
            setNavbarSettings(new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo));
        }

        if (secondHomeSection instanceof HTMLElement && heightToSecondSection <= HEIGHT_FROM_TOP_TO_FIRST_SECTION) {
            setNavbarSettings(new NavbarSettings('border-white', 'text-white', 'background-black', fourthLogo));
        }

        if (thirdHomeSection instanceof HTMLElement && heightToSecondSection <= HEIGHT_FROM_TOP_TO_THIRD_SECTION) {
            setNavbarSettings(new NavbarSettings('border-black', 'text-black', 'background-white', whiteLogo));
        }

        if (fourthHomeSection instanceof HTMLElement && heightToSecondSection <= HEIGHT_FROM_TOP_TO_FOURTH_SECTION) {
            setNavbarSettings(new NavbarSettings('border-white', 'text-white', 'background-pink', blackLogo));
        }
    };

    useEffect(() => {
        window.addEventListener('scroll', changeNavbarStyles);

        return () => window.removeEventListener('scroll', changeNavbarStyles);
    }, []);

    // The navbar links lead to the sections of this page: /#how-it-works and others.
    const { hash } = useLocation();
    useEffect(() => {
        if (hash) {
            document.getElementById(hash.replace(/^#/, ''))?.scrollIntoView({ behavior: 'smooth' });
        }
    }, [hash]);

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            <MintedBy />
            <HowItWorks />
            <FeaturedDrop />
            <ContactUs />
        </>
    );
};

export default Home;
