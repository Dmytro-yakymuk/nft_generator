import { useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import { RootState } from '@/app/store';
import { createHeader, getHeader, setHeader, updateHeaderImage, updateHeaderLogo, updateHeaderTextFields } from '@/app/store/actions/projects';
import { BadRequestError, NotFoundError } from '@/api';
import { convertToBase64 } from '@/internal/convertImage';

import UploadImage from '@components/common/UploadImage';

import './index.scss';

/** Initial index value. */
const FIRST_NODELIST_INDEX: number = 0;
const HEADER_LOGO_INDEX: number = 0;
const HEADER_BANNER_INDEX: number = 1;
const SECOND_INDEX: number = 1;
/** Website header links. */
const HEADER_NAV_LINKS: string[] = ['Collection', 'About', 'Roadmap', 'Team'];

const Header: React.FC = () => {
    const { id }: { id: string } = useParams();
    const dispatch = useDispatch();
    const header = useSelector((state: RootState) => state.collectionsReducer.mintsite.header);

    /** Changes header section content, i.e. title and description. */
    const changeHeaderTextContent = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name } = e.target;

        dispatch(updateHeaderTextFields({
            ...header, projectId: id,
            [name]: e.target.value,
        }));
    };

    useEffect(() => {
        (async() => {
            try {
                await dispatch(getHeader(id));
            } catch (error) {
                if (error instanceof NotFoundError) {
                    const title = '';
                    const description = '';
                    await dispatch(createHeader(id, title, description));
                }

                // TODO: error notification.
            }
        })();
    }, [id]);

    /** Event handler for adding header image and logo to header`s state. */
    const handleAddHeaderImages = async(e: React.ChangeEvent<HTMLInputElement>, index: number) => {
        if (!e.target || !e.target.files) {
            return;
        }

        /** Generate base64 string from image file. */
        const imageFile = e.target.files[FIRST_NODELIST_INDEX];
        const convertedImage: string = await convertToBase64(imageFile);
        const hashImageString = convertedImage.split(',')[SECOND_INDEX];

        if (index === HEADER_LOGO_INDEX) {
            await dispatch(setHeader({ ...header, projectId: id, logo: hashImageString }));
        }
        if (index === HEADER_BANNER_INDEX) {
            await dispatch(setHeader({ ...header, projectId: id, image: hashImageString }));
        }
    };

    /** Updates header fields. The project id is taken from the page address: the header in the store
     * keeps the default id until the section is loaded from the server. */
    const updateHeaderInformation = async() => {
        const projectHeader = { ...header, projectId: id };
        if (header.title || header.description) {
            await dispatch(updateHeaderTextFields(projectHeader));
        }

        header.logo && await dispatch(updateHeaderLogo(projectHeader));
        header.image && await dispatch(updateHeaderImage(projectHeader));
        location.reload();
    };

    return (
        <section className="header"
            onDrop={(e) => {
                e.preventDefault();
                e.stopPropagation();
            }}
            onDragOver={(e) => {
                e.preventDefault();
                e.stopPropagation();
            }}>
            <div className="header__wrapper">
                <div className="header-control">
                    <label htmlFor="header">Title</label>
                    <input
                        className="header-control__title"
                        type="text" id="header"
                        min="2"
                        max="150"
                        name="title"
                        value={header.title}
                        onChange={(e) => dispatch(setHeader({ ...header, title: e.target.value }))}
                    />
                    <UploadImage
                        title="Logo Image"
                        index={HEADER_LOGO_INDEX}
                        value={header.logo}
                        imageSize={'80x80'}
                        handleAddInputField={handleAddHeaderImages}
                    />
                    <UploadImage
                        title="Header Image"
                        index={HEADER_BANNER_INDEX}
                        value={header.image}
                        imageSize={'1440x164'}
                        handleAddInputField={handleAddHeaderImages}
                    />
                    <label className="header-control__label" htmlFor="description">Description</label>
                    <textarea
                        className="header-control__description"
                        id="description"
                        name="description"
                        value={header.description}
                        onChange={(e) => dispatch(setHeader({ ...header, description: e.target.value }))}
                    />
                    <button className="website-section-save" onClick={updateHeaderInformation}>
                        <span className="website-section-save__text">Save</span>
                    </button>
                </div>
                <div className="header-info">
                    <div className="header-info__wrapper">
                        <nav className="header-info__nav">
                            <div className="header-info__logo__wrapper">
                                <img className="header-info__logo" src={`${window.location.origin}/images/${id}/logo.png`} />
                            </div>
                            <ul className="header-info__links">
                                {
                                    HEADER_NAV_LINKS.map(
                                        (link: string, index: number) =>
                                            <li key={index} className="header-info__link">
                                                {link}
                                            </li>
                                    )
                                }
                            </ul>
                        </nav>
                        <div className="header-info__image__wrapper">
                            <img className="header-info__image" src={`${window.location.origin}/images/${id}/header.png`} />
                        </div>
                        <div className="header-info__text">
                            <span className="title">{header.title}</span>
                            <span className="description">{header.description}</span>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default Header;
