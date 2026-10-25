import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { ProjectItem } from '@components/common/Project';
import { Popup } from '@components/common/Popup';
import { Navbar } from '@components/common/Navbar';
import { NavbarSettings } from '@views/Home';

import pinkLogo from '@static/images/logo.png';

import { RootState } from '@/app/store';
import { InternalError } from '@/api';
import { createProject, getProjects } from '@/app/store/actions/projects';
import { Project } from '@/collections';
import { notify } from '@/internal/notifications';

import './index.scss';

const Projects = () => {
    /** Default project nft count. */
    const DEFAULT_NFT_COUNT: number = 1;
    const dispatch = useDispatch();
    const [isCreareProjectPopupVisible, setIsCreareProjectPopupVisible] = useState<boolean>(false);
    const [name, setName] = useState<string>('');
    const [nftCount, setNftCount] = useState<number>(DEFAULT_NFT_COUNT);
    const [projectNameWarning, setProjectNameWarning] = useState<string>('');
    const { projects } = useSelector((state: RootState) => state.collectionsReducer);
    const navbarSettings: NavbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);

    useEffect(() => {
        dispatch(getProjects());
    }, []);

    /** Opens create project popup. */
    const openCreateProjectPopup = () => {
        setProjectNameWarning('');
        setIsCreareProjectPopupVisible(true);
    };

    /** Changes project name. */
    const changeName = (e: React.ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value);
    };

    /** Changes project nft count. */
    const changeNftCount = (e: React.ChangeEvent<HTMLInputElement>) => {
        setNftCount(Number(e.target.value));
    };

    /** Adds and creates new project. */
    const addProject = async(e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setProjectNameWarning('');
        setNftCount(DEFAULT_NFT_COUNT);

        if (!name) {
            setProjectNameWarning('The field cannot be empty');

            return;
        }

        try {
            await dispatch(createProject(name, nftCount));
            setIsCreareProjectPopupVisible(false);
        } catch (error: any) {
            if (error instanceof InternalError) {
                setProjectNameWarning('This name is Already in use');

                return;
            }
            notify('Something went wrong, please try later');
        }
    };

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            {
                isCreareProjectPopupVisible &&
                <Popup
                    setIsPopupVisible={setIsCreareProjectPopupVisible}
                >
                    <h1 className="popup__title-left-side">
                        Create Project
                    </h1>
                    <form
                        className="create-project__form"
                        onSubmit={addProject}
                    >
                        <label
                            htmlFor="name"
                            className="create-project__form__name-label"
                        >
                            Project Name
                        </label>
                        <input
                            type="text"
                            className="create-project__form__name"
                            id="name"
                            min={1}
                            value={name}
                            onChange={changeName}
                        />
                        <span className="create-project__form__warning">{projectNameWarning}</span>
                        <label
                            htmlFor="nft"
                            className="create-project__form__nft-label"
                        >
                            Number Of NFTs
                        </label>
                        <input
                            type="number"
                            className="create-project__form__nft"
                            id="nft"
                            value={nftCount || ''}
                            onChange={changeNftCount}
                        />
                        <input
                            className="create-project__form__submit"
                            type="submit"
                            value="Create"
                        />
                    </form>
                </Popup>
            }
            <section className="projects">
                {projects.map((project: Project, index: number) =>
                    <ProjectItem
                        key={index}
                        project={project}
                    />
                )}
                <button
                    className="projects__create"
                    onClick={openCreateProjectPopup}
                >
                    New Project
                </button>
            </section>
        </>
    );
};

export default Projects;
