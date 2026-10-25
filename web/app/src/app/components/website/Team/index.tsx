import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import plus from '@static/images/collection-generator/plus.svg';
import trash from '@static/images/trash.svg';

import { RootState } from '@/app/store';
import { createMember, deleteMember, deleteMemberImage, getMembers, setMember, setMembers, updateMember, updateMemberImage } from '@/app/store/actions/projects';
import { convertToBase64 } from '@/internal/convertImage';
import { notify } from '@/internal/notifications';
import { Member } from '@/collections';

import './index.scss';

/** Number for increase/decrease paragraph value. */
const INCREASE_NUMBER: number = 1;

// TODO: should be reworked with children props.
const Team: React.FC = () => {
    const dispatch = useDispatch();
    const members = useSelector((state: RootState) => state.collectionsReducer.mintsite.team);
    const { id }: { id: string } = useParams();

    const projectId = id;

    /** Previewed images files. */
    const [previewedFiles, setPreviewedFiles] = useState<File[]>([]);

    /** Indicates if is first member image added. */
    const [isFirstMemberImageAdded, setIsFirstImageAdded] = useState<boolean>(false);

    /** Default members list length. */
    const DEFAULT_MEMBERS_LIST_LENGHT: number = 1;
    /** First nodelist element index. */
    const FIRST_NODELIST_INDEX: number = 0;
    /** Default members list length. */
    const FIRST_MEMBERS_LIST_ELEMENT_INDEX: number = 0;
    /** Default files list length. */
    const DEFAULT_FILES_LIST_LENGTH: number = 0;
    /** Converted image hash second index. */
    const CONVERTED_IMAGE_HASH_SECOND_INDEX: number = 1;

    /** Adds and saves new members to team. */
    const addMembers = async() => {
        const name = '';
        const role = '';
        try {
            if (members.length === DEFAULT_MEMBERS_LIST_LENGHT) {
                /** Saves default member. */
                members[FIRST_MEMBERS_LIST_ELEMENT_INDEX].id ?
                    await dispatch(updateMember(members[FIRST_MEMBERS_LIST_ELEMENT_INDEX])) :
                    await dispatch(createMember(projectId, members[FIRST_MEMBERS_LIST_ELEMENT_INDEX].name, members[FIRST_MEMBERS_LIST_ELEMENT_INDEX].role));
                /** Creates new member. */
                await dispatch(createMember(projectId, name, role));
                await dispatch(getMembers(projectId));
                setIsFirstImageAdded(true);

                return;
            }
            /** Updates and saves members if members are not saved before. */
            await Promise.all(
                members.map(
                    async(member: Member) => {
                        if (!member.isSaved) {
                            await dispatch(updateMember(member));
                        }
                    }
                )
            );
            await dispatch(createMember(projectId, name, role));
            await dispatch(getMembers(projectId));
        } catch (error) {
            notify('Members is not saved.');
        }
    };

    /** Saves members before publish website. */
    const saveMembers = async() => {
        try {
            if (members.length === DEFAULT_MEMBERS_LIST_LENGHT) {
                /** Saves default member. */
                members[FIRST_MEMBERS_LIST_ELEMENT_INDEX].id ?
                    await dispatch(updateMember(members[FIRST_MEMBERS_LIST_ELEMENT_INDEX])) :
                    await dispatch(createMember(projectId, members[FIRST_MEMBERS_LIST_ELEMENT_INDEX].name, members[FIRST_MEMBERS_LIST_ELEMENT_INDEX].role));
                /** Creates new member. */
                await dispatch(getMembers(projectId));
                setIsFirstImageAdded(true);

                return;
            }
            /** Updates and saves members if members are not saved before. */
            await Promise.all(
                members.map(
                    async(member: Member) => {
                        if (!member.isSaved) {
                            await dispatch(updateMember(member));
                        }
                    }
                )
            );
            await dispatch(getMembers(projectId));
            notify('Members saved!', 'success');
        } catch (error) {
            notify('Members is not saved.');
        }
    };

    /** Changes member image. */
    const changeMemberImage = (id: string) => async(e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files || e.target.files.length === DEFAULT_FILES_LIST_LENGTH) {
            return;
        }
        const copiedMember = members.find((member: Member) => member.id === id);
        if (!copiedMember) {
            return;
        }
        const index = members.findIndex((member: Member) => member.id === id);

        /** Sets images files to previewed files list. */
        const copiedPreviewedFiles = previewedFiles.slice();
        copiedPreviewedFiles[index] = e.target.files[FIRST_NODELIST_INDEX];
        setPreviewedFiles(copiedPreviewedFiles);

        const imageFile = e.target.files[FIRST_NODELIST_INDEX];
        /** Generate base64 string from image file. */
        const convertedImage: string = await convertToBase64(imageFile);
        const hashedImageString = convertedImage.split(',')[CONVERTED_IMAGE_HASH_SECOND_INDEX];

        await dispatch(updateMemberImage(projectId, id, hashedImageString));
        await dispatch(setMember(index, copiedMember));
    };

    /** Changes current member. */
    const changeMember = (id: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        let copiedMember = members.find((member: Member) => member.id === id);
        if (!copiedMember) {
            return;
        }
        const index = members.findIndex((member: Member) => member.id === id);
        copiedMember = {
            ...copiedMember,
            isSaved: false,
            [name]: value,
        };
        dispatch(setMember(index, copiedMember));
    };

    /** Creates a DOMString containing a URL representing the object given in the parameter. */
    const createImageURL = (image: File): string => URL.createObjectURL(image);

    /** Deletes current member. */
    const deleteCurrentMember = (memberId: string) => {
        dispatch(deleteMember(id, memberId));
        const filteredMembers = members.filter((member: Member) => member.id !== memberId);
        dispatch(setMembers(filteredMembers));
    };

    /** Deletes current member image on the server and its preview. */
    const deleteCurrentMemberImage = async(e: React.MouseEvent<HTMLButtonElement>, memberId: string, index: number) => {
        // the button sits inside the label of the file input: without this the file dialog opens too.
        e.preventDefault();
        await dispatch(deleteMemberImage(id, memberId));
        const copiedPreviewedFiles = previewedFiles.slice();
        delete copiedPreviewedFiles[index];
        setPreviewedFiles(copiedPreviewedFiles);
    };

    /** Returns actual member image. */
    const getPreviewImage = (previwedFileIndex: number, memberId: string): string =>
        previewedFiles[previwedFileIndex] ?
            createImageURL(previewedFiles[previwedFileIndex]) :
            `${window.location.origin}/images/${projectId}/members/${memberId}.png`;

    useEffect(() => {
        dispatch(getMembers(projectId));
    }, []);

    useEffect(() => {
        (async() => {
            const imageFile = previewedFiles[FIRST_NODELIST_INDEX];
            /** Generate base64 string from image file. */
            const convertedImage: string = await convertToBase64(imageFile);
            const hashedImageString = convertedImage.split(',')[CONVERTED_IMAGE_HASH_SECOND_INDEX];
            await dispatch(updateMemberImage(projectId, members[FIRST_MEMBERS_LIST_ELEMENT_INDEX].id, hashedImageString));
        })();
    }, [isFirstMemberImageAdded]);

    return (
        <section className="team"
            onDrop={(e) => {
                e.preventDefault();
                e.stopPropagation();
            }}
            onDragOver={(e) => {
                e.preventDefault();
                e.stopPropagation();
            }}>
            <div className="team__wrapper">
                <div className="team__control-area">
                    {
                        members.length &&
                        members.map(
                            (member: Member, index: number) =>
                                <div className="team__control-area__block" key={member.id}>
                                    <span className="team__control-area__title">Member {index + INCREASE_NUMBER}</span>
                                    <div
                                        className="team__control-area__btn"
                                        onClick={() => deleteCurrentMember(member.id)}
                                        id={`btn-${index}`}
                                    >
                                        <img src={trash} alt="Trash" />
                                        <button className="team__control-area__btn__delete">Delete</button>
                                    </div>
                                    <label htmlFor="name">
                                        Name
                                    </label>
                                    <input
                                        className="team__control-area__name"
                                        type="text"
                                        id="name"
                                        value={member.name}
                                        onChange={changeMember(member.id)}
                                        name="name"
                                    />
                                    <label htmlFor="role">
                                        Role (Optional)
                                    </label>
                                    <input
                                        className="team__control-area__role"
                                        type="text"
                                        id="role"
                                        value={member.role}
                                        onChange={changeMember(member.id)}
                                        name="role"
                                    />
                                    <div className="team__control-area__upload-image__wrapper">
                                        <label
                                            htmlFor={`upload-image-${index}`}
                                            id="upload-image-label"
                                        >
                                            {
                                                previewedFiles[index] ?
                                                    <div className="team__control-area__upload-image-selected">
                                                        <img
                                                            className="team__control-area__upload-image-selected__preview"
                                                            src={createImageURL(previewedFiles[index])}
                                                            alt="preview"
                                                        />
                                                        <span className="team__control-area__upload-image-selected__text">
                                                            {previewedFiles[index].name}
                                                        </span>
                                                        <button
                                                            className="team__control-area__upload-image-selected__btn"
                                                            onClick={(e) => deleteCurrentMemberImage(e, member.id, index)}
                                                        >
                                                            Delete
                                                        </button>
                                                    </div> : <span>
                                                        Drag and drop an image here, or click to browse
                                                    </span>
                                            }
                                        </label>
                                        <input
                                            className="team__control-area__upload-image"
                                            id={`upload-image-${index}`}
                                            accept="image/jpeg,image/png,image/gif"
                                            type="file"
                                            name="upload-image"
                                            hidden
                                            onChange={changeMemberImage(member.id)}
                                        />
                                    </div>
                                </div>
                        )
                    }
                    <div className="team__control-area__buttons">
                        <div className="team__control-area__add-member" onClick={addMembers}>
                            <img className="team__control-area__add-member__plus" src={plus} alt="Plus" />
                            <button className="team__control-area__add-member__btn">
                                <span className="team__control-area__add-member__btn__text">New Member</span>
                            </button>
                        </div>
                        <div className="team__control-area__save-member" onClick={saveMembers}>
                            <button className="team__control-area__save-member__btn">
                                <span className="team__control-area__save-member__btn__text">Save</span>
                            </button>
                        </div>
                    </div>
                </div>
                <div className="team__list">
                    <div className="team__list__wrapper">
                        <span className="team__list__title">Team</span>
                        <div className="team__list__members">
                            {
                                members.length &&
                                members.map(
                                    (member: Member, index: number) =>
                                        <div className="team__list__member" key={index}>
                                            <div className="team__list__member__logo__wrapper">
                                                <img
                                                    className="team__list__member__logo"
                                                    src={getPreviewImage(index, member.id)}
                                                    alt="Logo"
                                                />
                                            </div>
                                            <span className="team__list__member__name">{member.name}</span>
                                            <span className="team__list__member__type">{member.role}</span>
                                        </div>
                                )
                            }
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default Team;
