import { useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Route, Switch, useParams } from 'react-router-dom';

import { Navbar } from '@components/common/Navbar';
import { ProjectNavigation } from '@components/common/ProjectNavigation';
import { NavbarSettings } from '@views/Home';
import WebsiteNavigation from '@components/common/WebsiteNavigation';
import { Popup } from '@components/common/Popup';

import { RoutesConfig } from '@/app/routes';
import { RootState } from '@/app/store';
import { getProject, publishWebsite } from '@store/actions/projects';
import { BadRequestError } from '@/api';
import { notify } from '@/internal/notifications';

import pinkLogo from '@static/images/logo.png';

import './index.scss';

/** Website address: 3-63 lowercase letters, digits or hyphens that start and end with a letter or digit. */
const SLUG_PATTERN = /^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$/;
const SLUG_HINT = 'Use 3-63 lowercase letters, digits or hyphens, for example "my-collection"';

/** Suggests website address from the project name. */
const toSlug = (name: string): string =>
    name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '').slice(0, 63);

const Website: React.FC = () => {
    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();
    const project = useSelector((state: RootState) => state.collectionsReducer.project);
    const [domainName, setDomainName] = useState<string>('');
    const [inputError, setInputError] = useState<string>('');
    const [isPublishWebsiteVisible, setIsPublishWebsiteVisible] = useState<boolean>(false);
    const [isPublishedPopupVisible, setIsPublishedPopupVisible] = useState<boolean>(false);

    const domainInputRef = useRef<HTMLInputElement | null>(null);
    const navbarSettings: NavbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);
    const sitesPrefix = `${window.location.host}/sites/`;
    const publishedUrl = `${window.location.origin}/sites/${project.websiteSlug}`;
    // the same site at its own subdomain, if the server has a mint site domain.
    const mintSiteDomain = window.APP_CONFIG?.mintSiteDomain || '';
    const port = window.location.port ? `:${window.location.port}` : '';
    const subdomainUrl = mintSiteDomain && `${window.location.protocol}//${project.websiteSlug}.${mintSiteDomain}${port}/`;

    useEffect(() => {
        dispatch(getProject(id));
    }, [id]);

    useEffect(() => {
        if (project.id === id) {
            setDomainName(project.websiteSlug || toSlug(project.name));
        }
    }, [project.id, project.websiteSlug]);

    /** Publishes the website at the entered address. */
    const publish = async() => {
        setInputError('');

        if (!SLUG_PATTERN.test(domainName)) {
            setInputError(SLUG_HINT);

            return;
        }

        try {
            await dispatch(publishWebsite(id, domainName));
        } catch (error: any) {
            const message: string = error.message || '';
            if (message.includes('already taken')) {
                setInputError('This address is already taken');
            } else if (message.includes('not deployed')) {
                setInputError('Deploy the contract before publishing the website');
            } else if (message.includes('header section')) {
                setInputError('Fill in the Header section before publishing');
            } else if (error instanceof BadRequestError) {
                setInputError(SLUG_HINT);
            } else {
                notify('Could not publish the website');
            }

            return;
        }

        setIsPublishWebsiteVisible(false);
        setIsPublishedPopupVisible(true);
    };

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            <ProjectNavigation
                setIsPublishWebsiteVisible={setIsPublishWebsiteVisible}
                isPublishWebsiteVisible={isPublishWebsiteVisible}
            />
            <WebsiteNavigation />
            {
                isPublishWebsiteVisible &&
                    <Popup
                        setIsPopupVisible={setIsPublishWebsiteVisible}
                    >
                        <h1 className="popup__title-left-side">
                            Publish Website
                        </h1>
                        <div className="publish-project__domain">
                            <label htmlFor="domain">Address</label>
                            <div
                                className="publish-project__input"
                                onClick={() => {
                                    domainInputRef?.current?.focus();
                                }}
                                contentEditable
                                suppressContentEditableWarning={true}
                            >
                                <span className="domain">{sitesPrefix}</span>
                                <span className="project-name">{domainName}</span>
                            </div>
                            <input
                                type="text"
                                id="domain"
                                value={domainName}
                                onChange={(e) => setDomainName(e.target.value.toLowerCase())}
                                ref={domainInputRef}
                            />
                            <span className="publish-project__domain__warning">{inputError}</span>
                        </div>
                        <button className="publish-project__btn" onClick={publish}>
                            <span className="publish-project__btn__text">Publish</span>
                        </button>
                    </Popup>
            }
            {
                isPublishedPopupVisible &&
                    <Popup
                        setIsPopupVisible={setIsPublishedPopupVisible}
                    >
                        <h1 className="popup__title">
                            Your Website has been published
                        </h1>
                        <a
                            className="popup__open"
                            href={publishedUrl}
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            Open website
                        </a>
                        {subdomainUrl &&
                            <a
                                className="publish-project__subdomain"
                                href={subdomainUrl}
                                target="_blank"
                                rel="noopener noreferrer"
                            >
                                {subdomainUrl}
                            </a>
                        }
                    </Popup>
            }
            <Switch>
                {RoutesConfig.WebsiteEdit.children &&
                    RoutesConfig.WebsiteEdit.children.map(
                        (route, index: number) =>
                            <Route
                                key={index}
                                path={route.path}
                                component={route.component}
                                exact={route.exact}
                            />
                    )}
            </Switch>
        </>
    );
};

export default Website;
