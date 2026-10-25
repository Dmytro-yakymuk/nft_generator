import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';

import { Popup } from '@components/common/Popup';

import trash from '@static/images/project/trash.svg';

import { deleteProject } from '@store/actions/projects';
import { Project, ProjectStatus } from '@/collections';
import { notify } from '@/internal/notifications';
import { websiteTemplatePath } from '@/internal/websiteTemplatePath';

import './index.scss';

export const ProjectItem: React.FC<{ project: Project }> = ({ project }) => {
    const [isPopupVisible, setIsPopupVisible] = useState<boolean>(false);
    const [isReadyToPublishWebsite, setIsReadyToPublishWebsite] = useState<boolean>(false);
    const dispatch = useDispatch();
    const SUCCESS_STATUS_CODE = 200;
    const templateWebsiteUrl = websiteTemplatePath(project.id);

    /** Project title. */
    const title: string = `Delete "${project.name}"`;

    /** Deletes current project. */
    const deleteCurrentProject = async() => {
        try {
            await dispatch(deleteProject(project.id));
            setIsPopupVisible(!isPopupVisible);
        } catch (error) {
            notify('Something wrong!');
        }
    };

    /** Shows or closes popup. */
    const changePopupVisiblility = () => {
        setIsPopupVisible(!isPopupVisible);
    };

    /** Returns current project website status (published, draft or none status). */
    const currentProjectStatus = (isPublished: boolean): ProjectStatus => {
        if (isPublished) {
            return ProjectStatus.Published;
        }

        return isReadyToPublishWebsite ? ProjectStatus.Draft : ProjectStatus.None;
    };

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

    useEffect(() => {
        isPopupVisible ? document.body.classList.add('scroll-disabled')
            : document.body.classList.remove('scroll-disabled');

        return () => document.body.classList.remove('scroll-disabled');
    }, [isPopupVisible]);

    return <section className="project">
        <Link to={`/edit/${project.id}`}>
            <span className="project__title">Contract Name</span>
            <span className="project__name">{project.name}</span>
            <span className="project__title">Art</span>
            <span className="project__status">{project.areImagesGenerated ? ProjectStatus.Published : ProjectStatus.Draft}</span>
            <span className="project__title">Website</span>
            <span className="project__status">{currentProjectStatus(project.isWebsitePublished)}</span>
        </Link>
        <div className="project__footer">
            <span className="project__footer__text" onClick={changePopupVisiblility}>Delete</span>
            <img
                src={trash}
                alt="Trash"
                onClick={changePopupVisiblility}
            />
        </div>
        {isPopupVisible &&
            <Popup
                setIsPopupVisible={setIsPopupVisible}
            >
                <h1 className="popup__title-lower">
                    {title}
                </h1>
                <p className="popup__description">
                    Are you sure you want to delete this project?
                </p>
                <div className="popup__buttons">
                    <button
                        className="popup__confirm"
                        onClick={deleteCurrentProject}
                    >
                            Yes
                    </button>
                    <button
                        className="popup__cancel"
                        onClick={changePopupVisiblility}
                    >
                        No
                    </button>
                </div>
            </Popup>
        }
    </section>;
};
