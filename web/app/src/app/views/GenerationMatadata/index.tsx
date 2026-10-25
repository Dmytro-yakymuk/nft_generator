import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useParams } from 'react-router-dom';

import { Navbar } from '@components/common/Navbar';
import { ProjectNavigation } from '@components/common/ProjectNavigation';
import { NavbarSettings } from '@views/Home';

import { generateCollectionMetadata, getProject, pinIPFS, publishToIPFS } from '@/app/store/actions/projects';
import { InternalError } from '@/api';
import { RootState } from '@/app/store';
import { notify } from '@/internal/notifications';

import pinkLogo from '@static/images/logo.png';

import { ContractBaseURI } from './ContractBaseURI';
import { TokenGallery } from './TokenGallery';

import './index.scss';

const GenerationMatadata: React.FC = () => {
    const dispatch = useDispatch();
    const navbarSettings: NavbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);
    const project = useSelector((state: RootState) => state.collectionsReducer.project);
    const { id }: { id: string } = useParams();
    const history = useHistory();
    const [isPublishing, setIsPublishing] = useState(false);
    const [isPinning, setIsPinning] = useState(false);
    const gateway = window.APP_CONFIG?.ipfsGatewayUrl || 'https://ipfs.io';
    const pinningService = window.APP_CONFIG?.ipfsPinning || '';

    /** Generates collection metadata. */
    const generateMetadata = async() => {
        try {
            await generateCollectionMetadata(id);
            notify('Metadata generation was successful!', 'success');
        } catch (error: any) {
            if (error instanceof InternalError) {
                notify('Please generate the collection!');

                return;
            }
            notify('Could not generate metadata. Please, try again!');
        }
    };

    /** Opens the generator: the collection is generated again from there, after the layers or weights change. */
    const generateAgain = () => {
        history.push(`/collection-generator/${id}`);
    };

    /** Publishes the generated images and metadata to IPFS. */
    const publishCollectionToIPFS = async() => {
        setIsPublishing(true);
        try {
            await dispatch(publishToIPFS(id));
            notify('The collection has been published to IPFS!', 'success');
        } catch (error: any) {
            notify(`Could not publish the collection to IPFS: ${error.message}`);
        } finally {
            setIsPublishing(false);
        }
    };

    /** Pins the published images and metadata in the pinning service, so they stay in IPFS without the platform. */
    const pinCollection = async() => {
        setIsPinning(true);
        try {
            await dispatch(pinIPFS(id));
            notify(`The collection has been pinned in ${pinningService}!`, 'success');
        } catch (error: any) {
            notify(`Could not pin the collection: ${error.message}`);
        } finally {
            setIsPinning(false);
        }
    };

    useEffect(() => {
        dispatch(getProject(id));
    }, []);

    let publishButtonText = project.metadataCid ? 'Publish to IPFS Again' : 'Publish to IPFS';
    if (isPublishing) {
        publishButtonText = 'Publishing...';
    }

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            <ProjectNavigation />
            <section className="generate-metadata">
                <div className="generate-metadata__wrapper">
                    <h1 className="generate-metadata__title">Collection has been generated</h1>
                    <span className="generate-metadata__project-name">{project.name}</span>
                    <div className="generate-metadata__size">
                        <span className="label">Collection Size</span>
                        <span className="value">{project.nftCount}</span>
                    </div>
                    <div className="generate-metadata__dimensions">
                        <span className="label">Scene Dimensions</span>
                        <span className="value">500 x 500</span>
                    </div>
                    <button
                        className="generate-metadata__btn metadata"
                        onClick={generateMetadata}
                    >
                        <span className="generate-metadata__btn-text">Generate Metadata</span>
                    </button>
                    <div className="generate-metadata__divider"></div>
                    <button
                        className="generate-metadata__btn generate-again"
                        onClick={generateAgain}
                    >
                        <span className="generate-metadata__btn-text">Generate Again</span>
                    </button>
                    <button
                        className="generate-metadata__btn publish-to-ipfs"
                        onClick={publishCollectionToIPFS}
                        disabled={isPublishing}
                    >
                        <span className="generate-metadata__btn-text">{publishButtonText}</span>
                    </button>
                    {project.metadataCid &&
                        <div className="generate-metadata__ipfs">
                            <span className="label">Published to IPFS</span>
                            <a href={`${gateway}/ipfs/${project.imagesCid}/`} target="_blank" rel="noreferrer">
                                Images: {project.imagesCid}
                            </a>
                            <a href={`${gateway}/ipfs/${project.metadataCid}/0`} target="_blank" rel="noreferrer">
                                Metadata: {project.metadataCid}
                            </a>
                            <span className="value">Contract base URI: ipfs://{project.metadataCid}/</span>
                            {project.ipfsPinnedIn &&
                                <span className="value generate-metadata__pinned">
                                    Pinned in {project.ipfsPinnedIn}: available in IPFS without the platform
                                </span>
                            }
                            {pinningService && !project.ipfsPinnedIn &&
                                <button
                                    className="generate-metadata__btn pin-ipfs"
                                    onClick={pinCollection}
                                    disabled={isPinning}
                                >
                                    <span className="generate-metadata__btn-text">
                                        {isPinning ? 'Pinning...' : `Pin in ${pinningService}`}
                                    </span>
                                </button>
                            }
                            {project.id === id &&
                                <ContractBaseURI contractAddress={project.contractAddress} metadataCid={project.metadataCid} gateway={gateway} />
                            }
                        </div>
                    }
                </div>
                {project.id === id && <TokenGallery projectId={id} count={project.nftCount} />}
            </section>
        </>
    );
};

export default GenerationMatadata;
