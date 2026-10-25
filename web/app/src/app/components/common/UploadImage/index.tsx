import { useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { RootState } from '@/app/store';
import { setHeader } from '@/app/store/actions/projects';
import { Header } from '@/collections';
import { convertToBase64 } from '@/internal/convertImage';
import { notify } from '@/internal/notifications';

import './index.scss';

/** First choosen file index from input. */
const FIRST_FILE_INDEX: number = 0;
/** Initial drag counter number. */
const INITIAL_DRAG_COUNTER: number = 0;
/** Counter number. */
const COUNTER: number = 1;
/** Initial file length. */
const INITIAL_FILE_LENGTH: number = 0;

/** Warns when the image differs from the recommended size, e.g. 80x80: the site scales it, which may blur it. */
export const checkImageSize = (file: File | undefined, recommended: string) => {
    const [width, height] = recommended.split('x').map(Number);
    if (!file || !width || !height) {
        return;
    }
    const url = URL.createObjectURL(file);
    const image = new Image();
    image.onload = () => {
        if (image.naturalWidth !== width || image.naturalHeight !== height) {
            notify(`Recommended size is ${width}x${height} px, the image is ${image.naturalWidth}x${image.naturalHeight} px: it will be scaled.`, 'warning');
        }
        URL.revokeObjectURL(url);
    };
    image.src = url;
};

type UploadImageProps = {
    title: string;
    index?: number;
    value?: string;
    handleAddInputField?: (e: any, index: number) => void;
    imageSize: string;
};

const UploadImage: React.FC<UploadImageProps> = ({ title, index, value, handleAddInputField, imageSize }) => {
    const dispatch = useDispatch();
    const [isDragFile, setIsDragFile] = useState<boolean>(false);
    const [dragCounter, setDragCounter] = useState<number>(INITIAL_DRAG_COUNTER);
    const [uploadImage, setUploadImage] = useState<string>('');
    const header = useSelector((state: RootState) => state.collectionsReducer.mintsite.header);

    const uploadImageRef = useRef<HTMLDivElement | null>(null);
    const deleteImageBtnRef = useRef<HTMLButtonElement | null>(null);
    const imageInputRef = useRef<HTMLInputElement | null>(null);
    const imageLabelRef = useRef<HTMLLabelElement | null>(null);

    /** Handler for choosing or deleting an input image. */
    const imageInputHandler = () => {
        const input = imageInputRef.current;
        const label = imageLabelRef.current;

        if (!input) {
            return;
        }

        if (!label) {
            return;
        }

        /** Clearing the input when you click on the delete button. */
        if (deleteImageBtnRef.current) {
            input.value = '';
            label.innerText = 'Drag and drop an image here, or click to browse';

            setUploadImage('');
        }

        /** Writing the name and image source of the selected file. */
        input && label &&
            input.addEventListener('change', () => {
                input.files &&
                    (label.innerText = input.files[FIRST_FILE_INDEX].name);

                /* Create URL from input file. */
                const imageURL = URL?.createObjectURL(
                    // @ts-ignore
                    input.files && input.files[FIRST_FILE_INDEX]
                );

                setUploadImage(imageURL);
            });
    };

    /** The dragenter event is fired when a dragged element or text selection enters a valid drop target. */
    const handleDrag = (e: DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
    };

    /** The dragleave event is fired when a dragged element or text selection leaves a valid drop target. */
    const handleDragIn = (e: DragEvent) => {
        e.preventDefault();
        e.stopPropagation();

        setDragCounter((prev) => prev + COUNTER);

        if (
            e.dataTransfer &&
            e.dataTransfer.items &&
            e.dataTransfer.items.length > INITIAL_FILE_LENGTH
        ) {
            setIsDragFile(true);
        }
    };

    /** The dragover event is fired when an element or text selection is being dragged over a valid drop target. */
    const handleDragOut = (e: DragEvent) => {
        e.preventDefault();
        e.stopPropagation();

        setDragCounter((prev) => prev - COUNTER);

        if (dragCounter === INITIAL_DRAG_COUNTER) {
            setIsDragFile(false);
        }
    };

    /** Handles a file dropped on the block: it goes through the same input change handler as a chosen file. */
    const handleDrop = (e: DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setDragCounter(INITIAL_DRAG_COUNTER);
        setIsDragFile(false);

        const input = imageInputRef.current;
        const label = imageLabelRef.current;
        const file = e.dataTransfer?.files[FIRST_FILE_INDEX];
        if (!input || !label || !file) {
            return;
        }

        label.innerText = file.name;
        setUploadImage(URL.createObjectURL(file));

        // the section saves the image from the change event of the input, so the dropped file is put into
        // the input and the event is fired as if the file had been chosen in the dialog.
        const transfer = new DataTransfer();
        transfer.items.add(file);
        input.files = transfer.files;
        input.dispatchEvent(new Event('change', { bubbles: true }));
    };

    useEffect(() => {
        const uploadBlock = uploadImageRef?.current;

        uploadBlock?.addEventListener('dragenter', handleDragIn);
        uploadBlock?.addEventListener('dragleave', handleDragOut);
        uploadBlock?.addEventListener('dragover', handleDrag);
        uploadBlock?.addEventListener('drop', handleDrop);

        return () => {
            const uploadBlock = uploadImageRef?.current;

            uploadBlock?.removeEventListener('dragenter', handleDragIn);
            uploadBlock?.removeEventListener('dragleave', handleDragOut);
            uploadBlock?.removeEventListener('dragover', handleDrag);
            uploadBlock?.removeEventListener('drop', handleDrop);
        };
    }, []);

    /** htmlFor label, which depends of index. */
    const labelHtmlFor = index !== undefined ? index.toString() : 'image';
    /** Checking the index for truth (can be zero or null. */
    const isIndex = index || index === FIRST_FILE_INDEX;

    return (
        <div className="upload-image">
            <span className="upload-image__title">{title}</span>
            <div
                className={`image${uploadImage ? '-selected' : ''}`}
                onClick={() => imageInputHandler()}
                ref={uploadImageRef}
            >
                <img src={uploadImage.toString()} alt="image" />
                <label htmlFor={labelHtmlFor} id="choosen-img" ref={imageLabelRef}>
                Drag and drop an image here, or click to browse
                </label>
                <button ref={deleteImageBtnRef} className="image-selected__btn">Delete</button>
                {isIndex ?
                    <input
                        type="file"
                        accept="image/jpeg,image/png,image/gif"
                        id={index.toString()}
                        className="input-image"
                        hidden
                        ref={imageInputRef}
                        name="logo"
                        onChange={(e) => {
                            checkImageSize(e.target.files?.[FIRST_FILE_INDEX], imageSize);
                            handleAddInputField && handleAddInputField(e, index);
                        }}
                    /> :
                    <input
                        type="file"
                        accept="image/jpeg,image/png,image/gif"
                        id="image"
                        hidden
                        ref={imageInputRef}
                        name="logo"
                        onChange={() => {}}
                    />}
            </div>
            <span className="upload-image__notification">
                {`Recommended size: ${imageSize}px. JPG, PNG, or GIF. 10MB max
            size.`}
            </span>
        </div>
    );
};

export default UploadImage;
