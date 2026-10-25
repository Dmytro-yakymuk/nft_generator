import { useEffect, useState } from 'react';
import { DragDropContext, Draggable, DropResult, Droppable } from 'react-beautiful-dnd';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import { RootState } from '@/app/store';
import { deleteContentImage, getContentImageNames, updateContentImage } from '@/app/store/actions/projects';
import { convertToBase64 } from '@/internal/convertImage';
import { generateRandomId } from '@/internal/generateIdString';
import { notify } from '@/internal/notifications';

import trash from '@static/images/trash.svg';
import cross from '@static/images/collection-generator/cross.svg';

import './index.scss';

/** Default files list length. */
const DEFAULT_FILES_LIST_LENGTH: number = 0;
const SECOND_LIST_INDEX: number = 1;

const Content: React.FC = () => {
    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();
    const [previewedFiles, setPreviewedFiles] = useState<any[]>([]);
    const imageNames = useSelector((state: RootState) => state.collectionsReducer.mintsite.contentImageNames);

    // TODO: need to add  drug and drop.
    const handleDrag = (drag: DropResult) => {
        // drag.destination &&
        // dispatch(exchangeLayers(drag.source.index, drag.destination.index));
    };

    /** Changes selected files and set to local state state. */
    const setFiles = (e: React.ChangeEvent<HTMLInputElement>): void => {
        if (!e.target.files || e.target.files.length === DEFAULT_FILES_LIST_LENGTH) {
            return;
        }

        const copyPreviewedFiles = previewedFiles.slice();
        const files = Array.from(e.target.files);
        files.forEach(file => copyPreviewedFiles.push(file));
        setPreviewedFiles(copyPreviewedFiles);
    };

    /** Deletes image from previewed images. */
    const deleteFromPreviewedFiles = (deletedFile: File) => {
        const copyPreviewedFiles = previewedFiles.slice();
        const filteredFiels: string[] = copyPreviewedFiles.filter(file => file !== deletedFile);
        setPreviewedFiles(filteredFiels);
    };

    /** Create image Url from selected file. */
    const createImageUrl = (image: File): string => URL?.createObjectURL(image);

    /** Cets selected files to server. */
    const setContentImages = (): void => {
        previewedFiles.length &&
        previewedFiles.forEach(async(file: File) => {
            const convertedImage: string = await convertToBase64(file);
            const hashImageString = convertedImage.split(',')[SECOND_LIST_INDEX];
            const imageName = generateRandomId();

            try {
                await dispatch(updateContentImage(id, hashImageString, imageName));
                await dispatch(getContentImageNames(id));
                setPreviewedFiles([]);
            } catch (error) {
                notify('Something wrong! Try later!');
            }
            notify('Images saved!', 'success');
        });
    };

    useEffect(() => {
        dispatch(getContentImageNames(id));
    }, []);

    return (
        <section className="content">
            <div className="content__wrapper">
                <div className="content-control">
                    <span className="upload-image__title">Collection Images</span>
                    <div className="content-control__upload-image">
                        <label htmlFor="collection-content__upload">Drag and drop an image here, or click to browse</label>
                        <input
                            className="collection-content__upload-input"
                            id="collection-content__upload"
                            accept="image/jpeg,image/png,image/gif"
                            type="file"
                            name="file"
                            multiple
                            hidden
                            onChange={setFiles}
                        />
                    </div>
                    <span className="content-control__description">
                    Recommended size: 350x450px. JPG, PNG, or GIF. 10MB max size.
                    </span>
                    <div className="content-control__images">
                        <DragDropContext onDragEnd={handleDrag}>
                            <Droppable droppableId="dropArea1">
                                {(provided, snapshot) =>
                                    <div
                                        {...provided.droppableProps}
                                        ref={provided.innerRef}
                                    >
                                        {previewedFiles.map((image: File, index: number) =>
                                            <Draggable key={index} draggableId={`draggable${index}`} index={index}>
                                                {(provided, snapshot) =>
                                                    <div
                                                        ref={provided.innerRef}
                                                        {...provided.draggableProps}
                                                        {...provided.dragHandleProps}
                                                    >
                                                        <div className="image-block" key={index}>
                                                            <img className="image-block__image" src={createImageUrl(image)} alt="Image" />
                                                            <span className="image-block__name">{image.name}</span>
                                                            <img className="image-block__delete-btn"
                                                                src={trash}
                                                                alt="Trash"
                                                                onClick={() => deleteFromPreviewedFiles(image)}
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
                    <button className="website-section-save" onClick={setContentImages}>
                        <span className="website-section-save__text">Save</span>
                    </button>
                </div>
                <div className="content-info">
                    <div className="content-info__wrapper">
                        <div className="content-info__images">
                            {imageNames.map((name: string, index: number) =>
                                <div key={index}>
                                    <img
                                        className="content-info__delete-img"
                                        src={cross}
                                        alt="Delete"
                                        onClick={() => dispatch(deleteContentImage(id, name))}
                                    />
                                    <img

                                        className="content-info__logo"
                                        src={`${window.location.origin}/images/${id}/content/${name}`}
                                        alt="Image"
                                    />
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default Content;
