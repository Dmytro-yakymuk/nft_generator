import { Dispatch, SetStateAction, useState } from 'react';
import { useSelector } from 'react-redux';

import { RootState } from '@/app/store';
import { Layer } from '@/layers';
import { drawVariant, layerVariants } from '@/layers/rarity';

import './index.scss';

const PREVIEW_COUNT = 8;
const NO_IMAGES = 0;

/** Draws a variant of every layer by the saved weights: one image of the collection as it may come out. */
const drawImage = (layers: Layer[]) => layers
    .filter((layer: Layer) => layer.layerImagesCount > NO_IMAGES)
    .map((layer: Layer) => ({ layer: layer.name, variant: drawVariant(layerVariants(layer)) }));

/** Preview of random images composed of the layers by their saved weights, as the generation draws them. */
export const LayersPreview: React.FC<{ projectId: string; setIsPreviewVisible: Dispatch<SetStateAction<boolean>> }> = ({
    projectId,
    setIsPreviewVisible,
}) => {
    const layers = useSelector((state: RootState) => state.layersReducer.layers);
    const draw = () => Array.from({ length: PREVIEW_COUNT }, () => drawImage(layers));
    const [images, setImages] = useState(draw);

    return (
        <section className="layers-preview" onClick={() => setIsPreviewVisible(false)}>
            <div className="layers-preview__wrapper" onClick={(e) => e.stopPropagation()}>
                <div className="layers-preview__header">
                    <h2 className="layers-preview__title">Preview</h2>
                    <span className="layers-preview__hint">
                        Random images by the saved weights; the generation also makes every image unique.
                    </span>
                </div>
                <div className="layers-preview__grid">
                    {images.map((image, index) =>
                        <div className="layers-preview__image" key={index}
                            title={image.map(({ layer, variant }) => `${layer}: ${variant.name}`).join('\n')}>
                            {image.map(({ layer, variant }) => variant.image &&
                                <img
                                    key={layer}
                                    src={`${window.location.origin}/images/${projectId}/${layer}/${variant.image}`}
                                    alt={`${layer}: ${variant.name}`}
                                />
                            )}
                        </div>
                    )}
                </div>
                <div className="layers-preview__buttons">
                    <button className="layers-preview__shuffle" onClick={() => setImages(draw())}>Shuffle</button>
                    <button className="layers-preview__close" onClick={() => setIsPreviewVisible(false)}>Close</button>
                </div>
            </div>
        </section>
    );
};
