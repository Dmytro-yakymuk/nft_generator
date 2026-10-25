import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';

import { Navbar } from '@/app/components/common/Navbar';
import { Popup } from '@/app/components/common/Popup';
import { NavbarSettings } from '@views/Home';

import pinkLogo from '@static/images/logo.png';
import projectPreview from '@static/images/ProjectEdit/project.png';

import { deleteGeneratedNFTs, deleteMintsite, getProject } from '@/app/store/actions/projects';
import { getContract } from '@/app/store/actions/contracts';
import { NotFoundError } from '@/api';
import { RootState } from '@/app/store';
import { metamaskError, notify } from '@/internal/notifications';
import { ServicePlugin } from '@/plugins/service';
import { ProjectStatus } from '@/collections';

import './index.scss';

const ProjectEdit = () => {
    const { name } = useSelector((state: RootState) => state.collectionsReducer.project);
    /** Indicates if edit collection popup visible. */
    const [isEditCollectionPopupVisible, setIsEditCollectionPopupVisible] = useState<boolean>(false);
    /** Indicates if website popup visible. */
    const [isEditWebsitePopupVisible, setIsEditWebsitePopupVisivle] = useState<boolean>(false);
    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();
    const conractSignature = useSelector((state: RootState) => state.contractsReducer.contractSignature);
    const { areImagesGenerated, isWebsitePublished } = useSelector((state: RootState) => state.collectionsReducer.project);
    const ethersService = ServicePlugin.create();
    const navbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);
    // TODO: also this title could be changed on not empty string after design solution.
    /** Default popup title. */
    const DEFAULT_POPUP_TITLE: string = '';

    /** Changes edit collection popup visibility. */
    const changeEditCollectionPopupVisibility = () => {
        setIsEditCollectionPopupVisible(!isEditCollectionPopupVisible);
    };

    /** Changes contract popup visibility. */
    const changeContractPopupVisibility = () => {
        setIsEditWebsitePopupVisivle(!isEditWebsitePopupVisible);
    };

    /** Sets contract signature into reducer and deploys contract with this signature */
    const deployContract = async() => {
        await dispatch(getContract(id));
        try {
            await ethersService.deployContract(conractSignature);
        } catch (error: any) {
            metamaskError(error);
        }
    };

    /** Deletes project mintsite. */
    const deleteProjectMintSite = async() => {
        try {
            await deleteMintsite(id);
            setIsEditWebsitePopupVisivle(!isEditWebsitePopupVisible);
            notify('Website was delete.', 'success');
        } catch (error) {
            if (error instanceof NotFoundError) {
                notify('Website is not created.');

                return;
            }
            notify('Something wrong!');
        }
    };

    /** Deletes generates images, i.e. NFT's. */
    const deleteGeneratedImages = async() => {
        try {
            await deleteGeneratedNFTs(id);
            setIsEditCollectionPopupVisible(!isEditCollectionPopupVisible);
            notify('Collection was delete.', 'success');
        } catch (error) {
            notify('Failed to delete collection! Please try later!');
        }
    };

    useEffect(() => {
        dispatch(getProject(id));
    }, []);

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            <section className="project-edit">
                <div className="project-edit__navigation">
                    <h2 className="project-edit__navigation__heading">
                        {name}
                    </h2>
                    <ul className="project-edit__navigation__list">
                        <li className="project-edit__navigation__item" >
                            Art
                            <button
                                className="project-edit__navigation__item__btn"
                                onClick={changeEditCollectionPopupVisibility}
                            >
                                {areImagesGenerated ? ProjectStatus.Published : ProjectStatus.Draft}
                                <div className="project-edit__navigation__item__btn-draft" />
                            </button>
                        </li>
                        <li className="project-edit__navigation__item">
                            Website
                            <button
                                className="project-edit__navigation__item__btn"
                                onClick={changeContractPopupVisibility}
                            >
                                {isWebsitePublished ? ProjectStatus.Published : ProjectStatus.Draft}
                                <div className="project-edit__navigation__item__btn-draft" />
                            </button>
                        </li>
                    </ul>
                </div>
                <div className="project-edit__content">
                    <img
                        src={projectPreview}
                        alt="project preview"
                    />
                </div>
            </section>
            {
                isEditCollectionPopupVisible &&
                <Popup
                    setIsPopupVisible={setIsEditCollectionPopupVisible}
                >
                    {
                        <ul className="popup__list">
                            <li className="popup__list__item">
                                <Link
                                    to={`/collection-generator/${id}`}
                                    className="popup__list__item__link"
                                >
                                    {
                                        !areImagesGenerated ?
                                            'Edit collection' :
                                            'Open'
                                    }
                                </Link>
                            </li>
                            {
                                !areImagesGenerated &&
                                <li className="popup__list__item">
                                    <Link
                                        className="popup__list__item__link"
                                        to={`/create-contract/${id}`}
                                    >
                                        Edit contract
                                    </Link>
                                </li>
                            }
                            <li
                                className="popup__list__item"
                                onClick={deleteGeneratedImages}
                            >
                                Delete
                            </li>
                        </ul>
                    }
                </Popup>
            }
            {
                isEditWebsitePopupVisible &&
                <Popup
                    setIsPopupVisible={setIsEditWebsitePopupVisivle}
                >
                    {
                        <ul className="popup__list">
                            <li className="popup__list__item">
                                <Link
                                    to={`/website/${id}/header`}
                                    className="popup__list__item__link"
                                >
                                    {
                                        !isWebsitePublished ?
                                            'Edit' :
                                            'Open'
                                    }
                                </Link>
                            </li>
                            {isWebsitePublished &&
                                <li className="popup__list__item">
                                    <Link
                                        to={`/dashboard/${id}`}
                                        className="popup__list__item__link"
                                    >
                                        Dashboard
                                    </Link>
                                </li>
                            }
                            <li
                                className="popup__list__item"
                                onClick={deleteProjectMintSite}
                            >
                                Delete
                            </li>
                        </ul>
                    }
                </Popup>
            }
        </>
    );
};
export default ProjectEdit;
