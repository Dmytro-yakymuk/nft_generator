import { Navbar } from '@components/common/Navbar';
import { ProjectDashboard } from '@components/CollectionDashboard';

import pinkLogo from '@static/images/logo.png';

import { NavbarSettings } from '@views/Home';

const Dashboard = () => {
    const navbarSettings: NavbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            <ProjectDashboard />
        </>
    );
};

export default Dashboard;
