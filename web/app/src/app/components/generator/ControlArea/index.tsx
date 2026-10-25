import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useHistory, useParams } from 'react-router-dom';

import { GenerationProgress } from './GenerationProgress';
import { LayersList } from './LayersList';
import { Popup } from '@components/common/Popup';

import { RootState } from '@/app/store';
import { addLayer, getLayers } from '@/app/store/actions/layers';
import { generateNfts, getProject } from '@/app/store/actions/projects';
import { Layer } from '@/layers';
import { notify } from '@/internal/notifications';
import { BadRequestError } from '@/api';

import plus from '@static/images/collection-generator/plus.svg';

import './index.scss';

export const ControlArea: React.FC = () => {
    /** Delay after uploading images */
    const DELAY = 1000;
    const history = useHistory();
    const { id }: { id: string } = useParams();
    const projectId = id;
    const project = useSelector((state: RootState) => state.collectionsReducer.project);
    const layers = useSelector((state: RootState) => state.layersReducer.layers);
    const nftCount = useSelector((state: RootState) => state.collectionsReducer.nftGenerationCount);
    const dispatch = useDispatch();
    const [layerName, setLayerName] = useState<string>('');
    const [isPopUpVisible, setPopUpVisibility] = useState<boolean>(false);
    const [isGenerationPopUpVisible, setGenerationPopUpVisibility] = useState<boolean>(false);

    /** Creates new layer. */
    const createLayer = async () => {
        const isNameDuplicated = layers.find((layer: Layer) => layer.name === layerName);

        if (isNameDuplicated) {
            notify('Layer with such name already exists');

            return;
        }

        await dispatch(addLayer(projectId, layerName, layers.length));
        await setLayerName('');
        await setPopUpVisibility(false);
    };

    /** Changes popup visability. */
    const changePopupVisiblility = () => {
        setPopUpVisibility(!isPopUpVisible);
    };

    /** Starts NFT generation, and shows deneration progress popup */
    const handleNftGeneration = async () => {
        if (!layers.length) {
            notify('Please, create layer');

            return;
        }

        try {
            setGenerationPopUpVisibility(true);
            await dispatch(generateNfts(id));
        } catch (error: any) {
            notify(error.message.split(':').pop() || 'Something went wrong');
            setGenerationPopUpVisibility(false);
        }
    };

    useEffect(() => {
        dispatch(getLayers(projectId));
        dispatch(getProject(projectId));
    }, [projectId]);

    useEffect(() => {
        setTimeout(() => {
            if (project.nftCount && nftCount >= project.nftCount) {
                setGenerationPopUpVisibility(false);
                history.push(`/generate-metadata/${id}`);
            }
        }, DELAY);
    }, [nftCount]);

    return (
        <aside className="generator__aside">
            {isPopUpVisible &&
                <Popup
                    setIsPopupVisible={setPopUpVisibility}
                >
                    <h1 className="popup__title-lower">
                        Create a New Layer
                    </h1>
                    <input
                        type="text"
                        className="popup__change_name_input"
                        placeholder="Name"
                        value={layerName}
                        onChange={(e) => setLayerName(e.target.value)}
                    />
                    <div className="popup__buttons">
                        <button
                            className="popup__confirm"
                            onClick={createLayer}
                        >
                            Create
                        </button>
                        <button
                            className="popup__cancel"
                            onClick={changePopupVisiblility}
                        >
                            Cancel
                        </button>
                    </div>
                </Popup>
            }
            <div className="generator__aside__wrapper">
                <span className="generator__aside__project-name">{project.name}</span>
                <div className="generator__aside__collection-size">
                    <span className="title">Collection Size</span>
                    <span className="value">{project.nftCount}</span>
                </div>
                <div className="generator__aside__scene">
                    <span className="title">Scene Dimensions</span>
                    <span className="value">500 x 500</span>
                </div>
                <div className="generator__aside__new-layer" onClick={() => setPopUpVisibility(true)}>
                    <img className="generator__aside__new-layer__plus" src={plus} alt="Plus" />
                    <button className="generator__aside__new-layer__btn">
                        <span className="generator__aside__new-layer__btn__text">New Layer</span>
                    </button>
                </div>
                <LayersList />
            </div>
            <button
                className="generator__aside__confirm"
                onClick={handleNftGeneration}
            >
                Generate And Continue
            </button>
            {
                isGenerationPopUpVisible &&
                <GenerationProgress setGenerationPopUpVisibility={setGenerationPopUpVisibility} />
            }
        </aside>
    );
};
