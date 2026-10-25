import { useEffect, useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';

import { Popup } from '@components/common/Popup';

import art from '@static/images/project-navigation/art.svg';
import contract from '@static/images/project-navigation/contract.svg';
import website from '@static/images/project-navigation/website.svg';
import whitelist from '@static/images/project-navigation/whitelist.svg';

import { websiteTemplatePath } from '@/internal/websiteTemplatePath';

import './index.scss';

/** Projects navigation entity. */
class NavigationItem {
    public constructor(
        public name: string,
        public logo: string,
        public link: string,
    ) {}
}

type ProjectNavigationProps = {
    isPublishWebsiteVisible?: boolean;
    setIsPublishWebsiteVisible?: (isPopupVisible: boolean) => void;
};

export const ProjectNavigation: React.FC<ProjectNavigationProps> = ({ setIsPublishWebsiteVisible, isPublishWebsiteVisible }) => {
    const [isReadyToPublishWebsite, setIsReadyToPublishWebsite] = useState<boolean>(false);
    const location = useLocation();
    const SUCCESS_STATUS_CODE = 200;
    const regexp = /[a-z0-9\-]{36}/gi;
    const id = location.pathname.split('/').find(el => el.match(regexp));
    const publishButtonClassName: string = isReadyToPublishWebsite ? '' : '-disabled';
    const isCurrentPathExistWebsite = location.pathname.split('/').includes('website');
    const templateWebsiteUrl = websiteTemplatePath(id);

    /** Project navigation links. */
    const navigationItems: NavigationItem[] = [
        new NavigationItem('Art', art, `/collection-generator/${id}`),
        new NavigationItem('WhiteList', whitelist, `/whitelist/${id}`),
        new NavigationItem('Contract', contract, `/create-contract/${id}`),
    ];

    useEffect(() => {
        /** Checks website template on status code and changes local state. */
        fetch(templateWebsiteUrl)
            .then((response) =>
                response.status !== SUCCESS_STATUS_CODE
                    ?
                    setIsReadyToPublishWebsite(false)
                    :
                    setIsReadyToPublishWebsite(true));
    }, [isReadyToPublishWebsite]);

    return (
        <>
            <div className="project-navigation">
                <div className="project-navigation__wrapper">
                    <ul className="project-navigation__list">
                        {
                            isCurrentPathExistWebsite
                                ?
                                <NavLink to={`/website/${id}/about`} className="project-navigation__link" activeClassName="is-active">
                                    <img className="project-navigation__link__logo" src={website} alt="Website" />
                                    <span className="project-navigation__link__name">Website</span>
                                </NavLink>
                                :
                                navigationItems.map(
                                    (item: NavigationItem, index: number) =>
                                        <li key={index}>
                                            <NavLink to={item.link} className="project-navigation__link" activeClassName="is-active">
                                                <img className="project-navigation__link__logo" src={item.logo} alt={item.name} />
                                                <span className="project-navigation__link__name">{item.name}</span>
                                            </NavLink>
                                        </li>
                                )
                        }
                    </ul>
                    <a
                        // TODO: Need add logic to art link.
                        href={isCurrentPathExistWebsite ? templateWebsiteUrl : ''}
                        className={`project-navigation__template${isCurrentPathExistWebsite ? '' : '-hidden'}`}
                        target="_blank"
                    >
                        Save as draft
                    </a>
                    {
                        isCurrentPathExistWebsite &&
                        <button
                            className={`project-navigation__publish${publishButtonClassName}`}
                            onClick={() => setIsPublishWebsiteVisible && setIsPublishWebsiteVisible(true)}
                            disabled={!isReadyToPublishWebsite}
                        >
                            Publish
                        </button>
                    }
                </div>
            </div>
        </>
    );
};
