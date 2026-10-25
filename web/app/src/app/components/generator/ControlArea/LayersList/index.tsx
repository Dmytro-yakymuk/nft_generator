import { useState } from 'react';
import { DragDropContext, Draggable, DropResult, Droppable } from 'react-beautiful-dnd';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import { Popup } from '@/app/components/common/Popup';

import trash from '@static/images/trash.svg';

import { RootState } from '@/app/store';
import { deleteLayer, setActivelayer, updatePriority } from '@/app/store/actions/layers';
import { Layer } from '@/layers';

import './index.scss';

export const LayersList: React.FC = () => {
    /** TODO: replace after api changes */
    const DEFAULT_VALUE = 0;
    const dispatch = useDispatch();
    const { id }: {id: string} = useParams();
    const projectId = id;
    const layers = useSelector((state: RootState) => state.layersReducer.layers);
    const activeLayer = useSelector((state: RootState) => state.layersReducer.activeLayer);
    const [isPopUpVisible, setIsPopUpVisible] = useState(false);
    const [selectedLayerIndex, setSelectedLayerIndex] = useState(DEFAULT_VALUE);

    const handleDrag = (drag: DropResult) => {
        const selectedLayer = layers[drag.source.index];

        drag.destination &&
            dispatch(updatePriority(projectId, selectedLayer.id, drag.destination.index));
    };

    /** Select layer for delete. */
    const selectLayerForDelete = (e: any, index: number) => {
        e.stopPropagation();
        setSelectedLayerIndex(index);
        setIsPopUpVisible(!isPopUpVisible);
    };

    /** Changes popup visability. */
    const changePopupVisability = () => {
        setIsPopUpVisible(!isPopUpVisible);
    };

    /** Deletes current layer. */
    const deleteCurrentLayer = () => {
        dispatch(deleteLayer(layers[selectedLayerIndex].id, layers[selectedLayerIndex].projectId));
        setIsPopUpVisible(!isPopUpVisible);
    };

    return (
        <div className="generator__aside__layers">
            {isPopUpVisible &&
                <Popup
                    setIsPopupVisible={setIsPopUpVisible}
                >
                    <h1 className="popup__title-lower">
                        Delete "{layers[selectedLayerIndex].name}"
                    </h1>
                    <span className="popup__description">
                        Are you sure you want to delete this layer?
                    </span>
                    <div className="popup__buttons">
                        <button
                            className="popup__confirm"
                            onClick={deleteCurrentLayer}
                        >
                            Yes
                        </button>
                        <button
                            className="popup__cancel"
                            onClick={changePopupVisability}
                        >
                            No
                        </button>
                    </div>
                </Popup>
            }
            <DragDropContext onDragEnd={handleDrag} >
                <Droppable droppableId="dropArea1">
                    {(provided, snapshot) =>
                        <div
                            {...provided.droppableProps}
                            ref={provided.innerRef}
                        >
                            {layers.map((layer: Layer, index: number) =>
                                <Draggable draggableId={`draggable${index}`} index={index} key={index}>
                                    {(provided, snapshot) =>
                                        <div
                                            ref={provided.innerRef}
                                            {...provided.draggableProps}
                                            {...provided.dragHandleProps}
                                        >
                                            <div
                                                className={layer.id === activeLayer.id ? 'generator__aside__layer-active' : 'generator__aside__layer'}
                                                key={index}
                                                onClick={() => dispatch(setActivelayer(index))}
                                            >
                                                <div className="generator__aside__layer__count" >{layer.layerImagesCount}</div>
                                                <span className="generator__aside__layer__name">{layer.name}</span>
                                                {layer.noneWeight > DEFAULT_VALUE &&
                                                    <span className="generator__aside__layer__optional">optional</span>
                                                }
                                                <img
                                                    className="generator__aside__layer__delete"
                                                    src={trash} alt="Trash"
                                                    onClick={(e) => selectLayerForDelete(e, index)}
                                                />
                                            </div>
                                        </div>
                                    }
                                </Draggable>
                            )
                            }
                        </div>
                    }
                </Droppable>
            </DragDropContext>
        </div>
    );
};
