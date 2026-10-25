import { Navbar } from '@components/common/Navbar';
import { ControlArea } from '@components/generator/ControlArea';
import { CollectionContent } from '@/app/components/generator/CollectionContent';
import { ProjectNavigation } from '@components/common/ProjectNavigation';

import { NavbarSettings } from '@views/Home';

import pinkLogo from '@static/images/logo.png';

import './index.scss';

const CollectionGenerator: React.FC = () =>
    <>
        <Navbar navbarSettings={new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo)} />
        <ProjectNavigation />
        <section className="generator">
            <div className="generator__wrapper">
                <ControlArea />
                <CollectionContent />
            </div>
        </section>
    </>;


export default CollectionGenerator;
