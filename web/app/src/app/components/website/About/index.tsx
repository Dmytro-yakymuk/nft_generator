import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import UploadImage from '@components/common/UploadImage';

import { NotFoundError } from '@/api';
import { RootState } from '@/app/store';
import { createAboutSection, getAboutSection, setAboutSection, updateAboutSection } from '@/app/store/actions/projects';
import { convertToBase64 } from '@/internal/convertImage';
import { AboutSection } from '@/collections';

import './index.scss';

const About: React.FC = () => {
    const aboutSection = useSelector((state: RootState) => state.collectionsReducer.mintsite.aboutSection);
    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();

    const projectId = id;

    /** First nodelist element index. */
    const FIRST_NODELIST_INDEX: number = 0;
    /** About section image index. */
    const ABOUT_SECTION_IMAGE_INDEX: number = 0;
    /** Converted image hash second index. */
    const CONVERTED_IMAGE_HASH_SECOND_INDEX: number = 1;

    /** Image size placeholder. */
    const IMAGE_SIZE_PLACEHOLDER: string = '240x300';

    const aboutSectionImagePath = `${window.location.origin}/images/${projectId}/about.png`;

    /** Changes about section content, i.e. title and description. */
    const changeAboutSectionContent = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name } = e.target;

        dispatch(setAboutSection({
            ...aboutSection,
            [name]: e.target.value,
        }));
    };

    /** Changes 'about' section image. */
    const chnageAboutSectionImage = async(e: React.ChangeEvent<HTMLInputElement>, index: number) => {
        if (!e.target || !e.target.files) {
            return;
        }

        const imageFile = e.target.files[FIRST_NODELIST_INDEX];
        /** Generate base64 string from image file. */
        const convertedImage: string = await convertToBase64(imageFile);
        const hashedImageString = convertedImage.split(',')[CONVERTED_IMAGE_HASH_SECOND_INDEX];

        await dispatch(setAboutSection({
            ...aboutSection,
            image: hashedImageString,
        }));
    };

    /** Updates ( i.e saves) about section. */
    const saveAboutSection = async() => {
        await dispatch(updateAboutSection({ ...aboutSection, projectId }));
        location.reload();
    };

    useEffect(() => {
        (async() => {
            try {
                await dispatch(getAboutSection(projectId));
            } catch (error) {
                if (error instanceof NotFoundError) {
                    const title = '';
                    const description = '';
                    await dispatch(createAboutSection(id, title, description));
                }
            }
        })();
    }, []);

    return <section className="about"
        onDrop={(e) => {
            e.preventDefault();
            e.stopPropagation();
        }}
        onDragOver={(e) => {
            e.preventDefault();
            e.stopPropagation();
        }}>
        <div className="about__wrapper">
            <div className="about-control">
                <label htmlFor="about">Title</label>
                <input
                    className="about-control__title"
                    type="text"
                    id="about"
                    placeholder="About"
                    name="title"
                    value={aboutSection.title}
                    onChange={changeAboutSectionContent}
                />
                <UploadImage
                    title="Image"
                    value={aboutSection.image}
                    imageSize={IMAGE_SIZE_PLACEHOLDER}
                    index={ABOUT_SECTION_IMAGE_INDEX}
                    handleAddInputField={chnageAboutSectionImage}
                />
                <label className="about-control__label" htmlFor="description">Description</label>
                <textarea
                    className="about-control__description"
                    id="description"
                    name="description"
                    value={aboutSection.description}
                    onChange={changeAboutSectionContent}
                />
                <button
                    className="website-section-save"
                    onClick={saveAboutSection}
                >
                    <span className="website-section-save__text">
                        Save
                    </span>
                </button>
            </div>
            <div className="about-info">
                <div className="about-info__wrapper">
                    <div className="about-info__text">
                        <span className="title">
                            {aboutSection.title}
                        </span>
                        <span className="description">
                            {aboutSection.description}
                        </span>
                    </div>
                    <div className="about-info__logo__wrapper">
                        <img
                            className="about-info__logo"
                            src={aboutSectionImagePath}
                            alt="Logo"
                        />
                    </div>
                </div>
            </div>
        </div>
    </section>;
};
export default About;
