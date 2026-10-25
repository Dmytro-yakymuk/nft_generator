import { NavLink, useParams } from 'react-router-dom';

import './index.scss';

/** Exposes website navigation entity. */
class Navigation {
    /** At least, contains label and path. */
    constructor(
        public label: string,
        public path: string,
    ) { };
};

const WebsiteNavigation: React.FC = () => {
    const { id }: { id: string } = useParams();

    /** Website navigation config. */
    const navigationConfig: Navigation[] = [
        new Navigation('About', `/website/${id}/about`),
        new Navigation('Header', `/website/${id}/header`),
        new Navigation('Content', `/website/${id}/content`),
        new Navigation('Team', `/website/${id}/team`),
        new Navigation('Roadmap', `/website/${id}/roadmap`),
        new Navigation('Footer', `/website/${id}/footer`),
    ];

    return (
        <div className="website-navigation">
            <div className="website-navigation__wrapper">
                <h1 className="website-navigation__title">Website</h1>
                <div className="website-navigation__links">
                    {navigationConfig.map((navigation, index) =>
                        <NavLink key={index} to={navigation.path} activeClassName="is-active">
                            <span className="link">{navigation.label}</span>
                        </NavLink>
                    )}
                </div>
            </div>
        </div>
    );
};

export default WebsiteNavigation;
