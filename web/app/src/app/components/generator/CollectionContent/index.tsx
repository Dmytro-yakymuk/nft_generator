import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import { RootState } from '@/app/store';
import { deleteImage, getFilesNames, updateImageName, updateWeights, uploadImage } from '@/app/store/actions/layers';
import { ImagesNameIdentificators, LayerImagesIdentificators } from '@/layers';
import { DEFAULT_WEIGHT, MAX_WEIGHT, NONE_NAME, percentOf, rarityTier } from '@/layers/rarity';
import '@/layers/rarity.scss';
import { LayersPreview } from '@/app/components/generator/LayersPreview';
import { FileDimensions } from '@/app/types/generator';
import { Popup } from '@components/common/Popup';
import { convertToBase64 } from '@/internal/convertImage';
import { notify } from '@/internal/notifications';

import cross from '@static/images/collection-generator/cross.svg';
import pensil from '@static/images/collection-generator/pensil.svg';
import noImages from '@static/images/collection-generator/noImages.svg';

import './index.scss';

export const CollectionContent = () => {
    const DEFAULT_VALUE = 0;
    const DIMENSION = 500;
    /** Weight of being absent of a mandatory layer and the least weight of an image. */
    const MANDATORY = 0;
    const MIN_WEIGHT = 1;
    const PERCENT_DIGITS = 1;
    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();
    const projectId = id;
    /** Contains selected active layer. */
    const activeLayer = useSelector((state: RootState) => state.layersReducer.activeLayer);
    const layers = useSelector((state: RootState) => state.layersReducer.layers);
    const fileNames = useSelector((state: RootState) => state.layersReducer.fileNames);
    /** Contains files dimesnions */
    const [dimensions, setDimensions] = useState<FileDimensions[]>([]);
    /** Contains active card index */
    const [activeImageIndex, setActiveImageIndex] = useState<number>(DEFAULT_VALUE);
    /** Indicates if rename image popup is visible. */
    const [isRenameImagePopupVisible, setIsRenameImagePopupVisible] = useState<boolean>(false);
    /** Contains changed name for active image */
    const [name, setName] = useState('');
    /** Weights being edited, by image file name, as typed. */
    const [draftWeights, setDraftWeights] = useState<Record<string, string>>({});
    /** Weight of the layer being absent being edited; 0 makes the layer mandatory. */
    const [draftNoneWeight, setDraftNoneWeight] = useState<string>('0');
    /** Indicates if the preview of random images is visible. */
    const [isPreviewVisible, setIsPreviewVisible] = useState<boolean>(false);

    /** The edited weights start from the saved ones whenever the layer or its weights change. */
    useEffect(() => {
        const weights: Record<string, string> = {};
        Object.keys(activeLayer.weights || {}).forEach((image) => {
            weights[image] = String(activeLayer.weights[image]);
        });
        setDraftWeights(weights);
        setDraftNoneWeight(String(activeLayer.noneWeight || MANDATORY));
    }, [activeLayer]);

    /** Returns the typed weight as a number, or NaN if it is not an allowed weight. */
    const parseWeight = (value: string | undefined, min: number): number => {
        const weight = Number(value === undefined ? DEFAULT_WEIGHT : value);

        return Number.isInteger(weight) && weight >= min && weight <= MAX_WEIGHT ? weight : NaN;
    };
    const isOptional = draftNoneWeight !== String(MANDATORY);
    const imageWeights = fileNames.map((image: string) => parseWeight(draftWeights[image], MIN_WEIGHT));
    const noneWeight = isOptional ? parseWeight(draftNoneWeight, MIN_WEIGHT) : MANDATORY;
    const areWeightsValid = imageWeights.every((weight: number) => !isNaN(weight)) && !isNaN(noneWeight);
    const totalWeight = areWeightsValid ? imageWeights.reduce((sum: number, weight: number) => sum + weight, noneWeight) : MANDATORY;
    const areWeightsChanged = fileNames.some((image: string, index: number) =>
        imageWeights[index] !== (activeLayer.weights && activeLayer.weights[image] || DEFAULT_WEIGHT))
        || noneWeight !== (activeLayer.noneWeight || MANDATORY);

    /** Saves the edited weights of the active layer. */
    const saveWeights = async() => {
        const weights: Record<string, number> = {};
        fileNames.forEach((image: string, index: number) => {
            weights[image] = imageWeights[index];
        });
        try {
            await dispatch(updateWeights(projectId, activeLayer.id, noneWeight, weights));
            notify('Weights saved');
        } catch (error: any) {
            notify(error.message.split(':').pop() || 'Could not save weights');
        }
    };

    /** Makes the layer optional with a default weight of being absent, or mandatory. */
    const toggleOptional = () => {
        setDraftNoneWeight(String(isOptional ? MANDATORY : DEFAULT_WEIGHT));
    };

    /** Shows the weight, share and rarity tier of a variant. */
    const weightControl = (value: string | undefined, weight: number, onChange: (value: string) => void) => {
        const percent = percentOf(weight, totalWeight);
        const tier = rarityTier(percent);

        return <div className="collection-image__weight">
            <input
                type="number"
                min={1}
                max={MAX_WEIGHT}
                className={`collection-image__weight__input${isNaN(weight) ? '-error' : ''}`}
                value={value === undefined ? String(DEFAULT_WEIGHT) : value}
                onChange={(e) => onChange(e.target.value)}
                aria-label="weight"
            />
            {areWeightsValid &&
                <>
                    <span className="collection-image__weight__percent">{percent.toFixed(PERCENT_DIGITS)}%</span>
                    <span className={`rarity rarity--${tier.toLowerCase()}`}>{tier}</span>
                </>
            }
        </div>;
    };

    /** Default files list length. */
    const DEFAULT_FILES_LIST_LENGTH: number = 0;
    const SECOND_INDEX: number = 1;

    /** Returns loaded image dimensions */
    const handleDimensions = (name: string) => new Promise<FileDimensions>((resolve, reject) => {
        const image = new Image();
        image.src = `${window.location.origin}/images/${projectId}/${activeLayer.name}/${name}`;

        image.addEventListener('load', () => {
            resolve(new FileDimensions(image.naturalHeight, image.naturalWidth));
        });
    });
    /** Creates a preview as a side effect, whenever selected files is changed. */
    useEffect(() => {
        if (!fileNames) {
            return;
        };

        async function createDimensions() {
            const dimensions = await Promise.all(fileNames.map(async(name: string) => await handleDimensions(name))).then((result: FileDimensions[]) => result);
            setDimensions(dimensions);
        }
        createDimensions();
    }, [fileNames]);


    /** Changes selected files. */
    const uploadFiles = async(e: React.ChangeEvent<HTMLInputElement>) => {
        if (!layers.length) {
            notify('Please, create layer');

            return;
        }

        if (!e.target.files || e.target.files.length === DEFAULT_FILES_LIST_LENGTH) {
            return;
        }

        const files = Array.from(e.target.files);
        files.forEach(async(file) => {
            const convertedFile: string = await convertToBase64(file);
            const layerImagesParameters = new LayerImagesIdentificators(
                projectId, activeLayer.id, file.name, convertedFile.split(',')[SECOND_INDEX]
            );
            await dispatch(uploadImage(layerImagesParameters));
        });

        await dispatch(getFilesNames(projectId, activeLayer.id));
    };

    /** Deletes selected image. */
    const deleteSelectedImage = (imageName: string) => {
        dispatch(deleteImage(projectId, activeLayer.id, imageName));
    };

    /** Handles drag event and sends files to the api */
    const handleImagesDrag = async(e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();

        if (!layers.length) {
            notify('Please, create layer');

            return;
        }

        if (!e.dataTransfer.files || e.dataTransfer.files.length === DEFAULT_FILES_LIST_LENGTH) {
            return;
        }

        const files: File[] = Array.from(e.dataTransfer.files);
        files.forEach(async(file) => {
            const convertedFile: string = await convertToBase64(file);
            const layerImagesParameters = new LayerImagesIdentificators(
                projectId, activeLayer.id, file.name, convertedFile.split(',')[SECOND_INDEX]
            );
            await dispatch(uploadImage(layerImagesParameters));
        });
        await dispatch(getFilesNames(projectId, activeLayer.id));
    };

    useEffect(() => {
        // until a layer is chosen there is nothing to list: /layers/ without an id is not an API route.
        if (activeLayer.id) {
            dispatch(getFilesNames(projectId, activeLayer.id));
        }
    }, [activeLayer]);

    /** Checks image dimensions and returns boolean */
    function isDimensionsAppropriate(sizes: FileDimensions) {
        return sizes.height === DIMENSION || sizes.width === DIMENSION;
    }

    /** RenamesImageName */
    const renameImageName = () => {
        const imageParameters = new ImagesNameIdentificators(projectId, activeLayer.id, fileNames[activeImageIndex], name);
        dispatch(updateImageName(imageParameters));
        dispatch(getFilesNames(projectId, activeLayer.id));
        setIsRenameImagePopupVisible(!isRenameImagePopupVisible);
    };

    /** Canceles rename image. */
    const cancelRenameImage = () => {
        setIsRenameImagePopupVisible(!isRenameImagePopupVisible);
    };

    /** Shows popup for name editing and sets clicked input as active into state. */
    const editName = (index: number) => {
        setIsRenameImagePopupVisible(!isRenameImagePopupVisible);
        setActiveImageIndex(index);
    };

    return <section className="collection-content">
        {
            isRenameImagePopupVisible &&
            <Popup
                setIsPopupVisible={setIsRenameImagePopupVisible}
            >
                <h1 className="popup__title-lower">
                    Rename Image
                </h1>
                <input
                    type="text"
                    className="popup__change_name_input"
                    placeholder="Name"
                    onChange={(e) => setName(e.target.value)}
                />
                <div className="popup__buttons">
                    <button
                        className="popup__confirm"
                        onClick={renameImageName}
                    >
                        Confirm
                    </button>
                    <button
                        className="popup__cancel"
                        onClick={cancelRenameImage}
                    >
                        Cancel
                    </button>
                </div>
            </Popup>
        }
        {isPreviewVisible && <LayersPreview projectId={projectId} setIsPreviewVisible={setIsPreviewVisible} />}
        <div className="collection-content__header">
            <h2 className="collection-content__heading">{activeLayer.name}</h2>
            {layers.length > DEFAULT_FILES_LIST_LENGTH &&
                <div className="collection-content__controls">
                    <label className="collection-content__optional">
                        <input type="checkbox" checked={isOptional} onChange={toggleOptional} />
                        Optional layer
                    </label>
                    <button className="collection-content__preview" onClick={() => setIsPreviewVisible(true)}>
                        Preview
                    </button>
                </div>
            }
        </div>
        <div className="collection-content__images-area"
            onDrop={(e: React.DragEvent<HTMLDivElement>) => handleImagesDrag(e)}
            onDragOver={(e) => {
                e.preventDefault();
                e.stopPropagation();
            }}>
            {isOptional && fileNames.length > DEFAULT_FILES_LIST_LENGTH &&
                <div className="collection-image">
                    <div className="collection-image__wrapper collection-image__none">
                        <span>No {activeLayer.name}</span>
                    </div>
                    <div className="collection-image__text-wrapper">
                        <p className="collection-image__name">{NONE_NAME}</p>
                        {weightControl(draftNoneWeight, noneWeight, setDraftNoneWeight)}
                    </div>
                </div>
            }
            {fileNames.length ?
                fileNames.map((name: string, index: number) =>
                    <div className="collection-image" key={index}>
                        <div className="collection-image__wrapper">
                            <img className="collection-image__delete"
                                src={cross} alt="delete button"
                                onClick={() => deleteSelectedImage(fileNames[index])} />
                            <img
                                src={`${window.location.origin}/images/${projectId}/${activeLayer.name}/${name}`}
                                className="collection-image__file"
                                alt="layer media"
                            />
                        </div>
                        <div className="collection-image__text-wrapper">
                            <p className="collection-image__name">{name}</p>
                            {
                                dimensions[index] &&
                                <p className={`collection-image__image-dimensions${isDimensionsAppropriate(dimensions[index]) ? '' : '-error'}`}>
                                    {
                                        `${dimensions[index].height} X ${dimensions[index].width}`
                                    }
                                </p>
                            }
                            <img className="collection-image__edit"
                                src={pensil} alt="delete button"
                                onClick={() => editName(index)} />
                            {weightControl(draftWeights[name], imageWeights[index],
                                (value: string) => setDraftWeights({ ...draftWeights, [name]: value }))}
                        </div>
                    </div>
                )
                :
                <div className="collection-empty">
                    <img src={noImages} alt="no images" />
                    <p className="collection-empty__title">
                        No Images
                    </p>
                    <span className="collection-empty__description">
                        Please drag them here or click Upload
                    </span>
                </div>
            }
        </div>
        <div className="collection-content__wrapper">
            <div className="collection-content__description">
                Images must be 500x500 PNG. A weight sets how often the image appears:
                its share is the weight divided by the sum of the layer weights.
            </div>
            <button
                className="collection-content__save"
                disabled={!areWeightsChanged || !areWeightsValid || !fileNames.length}
                onClick={saveWeights}
            >
                Save Weights
            </button>
            <label htmlFor="collection-content__upload">Upload</label>
            <input
                className="collection-content__upload"
                id="collection-content__upload"
                type="file"
                name="file"
                accept=".png"
                multiple
                hidden
                onChange={uploadFiles}
            />
        </div>
    </section >;
};
