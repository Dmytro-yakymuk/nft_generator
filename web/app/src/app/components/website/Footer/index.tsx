import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import plus from '@static/images/collection-generator/plus.svg';

import { RootState } from '@/app/store';
import { createSocialNetwork, deleteSocialNetwork, getSocialNetworks, setSocialNetwork, setSocialNetworks, updateSocialNetwork } from '@/app/store/actions/projects';
import { SocialNetwork } from '@/collections';
import { notify } from '@/internal/notifications';

import './index.scss';

/** Number for increase paragraph value. */
const INCREASE_NUMBER: number = 1;

// TODO: should be reworked with children props.
const Footer: React.FC = () => {
    /** Default social networks list length. */
    const DEFAULT_SOCIAL_NETWORKS_LIST_LENGTH: number = 1;
    /** Default social networks list length. */
    const FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX: number = 0;

    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();
    const socialNetworks = useSelector((state: RootState) => state.collectionsReducer.mintsite.socialNetworks);

    /** Change social network. */
    const changeSocialNetwork = (id: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        let copiedSocialNetwork = socialNetworks.find((socialNetwork: SocialNetwork) => socialNetwork.id === id);

        if (!copiedSocialNetwork) {
            return;
        }

        const index = socialNetworks.findIndex((socialNetwork: SocialNetwork) => socialNetwork.id === id);
        copiedSocialNetwork = { ...copiedSocialNetwork, [name]: value };
        dispatch(setSocialNetwork(index, copiedSocialNetwork));
    };

    /** Adds new social network. */
    const addSocialNetwork = async() => {
        const title = '';
        const link = '';
        try {
            if (socialNetworks.length === DEFAULT_SOCIAL_NETWORKS_LIST_LENGTH) {
                socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX].id ?
                    await dispatch(updateSocialNetwork(socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX])) :
                    await dispatch(createSocialNetwork(id, socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX].title, socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX].link));
                await dispatch(createSocialNetwork(id, title, link));
                await dispatch(getSocialNetworks(id));

                return;
            }
            await Promise.all(
                socialNetworks.map(
                    async(socialNetwork: SocialNetwork) => await dispatch(updateSocialNetwork(socialNetwork))
                )
            );
            await dispatch(createSocialNetwork(id, title, link));
            await dispatch(getSocialNetworks(id));
            notify('Footer link saved!', 'success');
        } catch (error) {
            notify('Footer is not saved!');
        }
    };

    /** Saves new social network. */
    const saveSocialNetworks = async() => {
        try {
            if (socialNetworks.length === DEFAULT_SOCIAL_NETWORKS_LIST_LENGTH) {
                socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX].id ?
                    await dispatch(updateSocialNetwork(socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX])) :
                    await dispatch(createSocialNetwork(id, socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX].title, socialNetworks[FIRST_SOCIAL_NETWORKS_LIST_ELEMENT_INDEX].link));
                await dispatch(getSocialNetworks(id));

                return;
            }
            await Promise.all(
                socialNetworks.map(
                    async(socialNetwork: SocialNetwork) => await dispatch(updateSocialNetwork(socialNetwork))
                )
            );
            await dispatch(getSocialNetworks(id));
            notify('Footer saved!', 'success');
        } catch (error) {
            notify('Footer is not saved!');
        }
    };

    /** Deletes current social network. */
    const deleteCurrentSocialNetwork = async(socialNetworkId: string) => {
        await dispatch(deleteSocialNetwork(id, socialNetworkId));
        const filteredSocialNetworks = socialNetworks.filter((socialNetwork: SocialNetwork) => socialNetwork.id !== socialNetworkId);
        await dispatch(setSocialNetworks(filteredSocialNetworks));
    };

    useEffect(() => {
        dispatch(getSocialNetworks(id));
    }, []);

    return (
        <section className="footer">
            <div className="footer__wrapper">
                <div className="footer__control-area">
                    {
                        socialNetworks.map(
                            (socialNetwork: SocialNetwork, index: number) =>
                                <div className="footer__control-area__block" key={socialNetwork.id}>
                                    <span className="footer__control-area__title">Paragraph {index + INCREASE_NUMBER}</span>
                                    <a
                                        onClick={() => deleteCurrentSocialNetwork(socialNetwork.id)}
                                        className="footer__control-area__delete"
                                        id={`btn-${index}`}
                                    >
                                        Delete
                                    </a>
                                    <label htmlFor="title">Title</label>
                                    <input
                                        className="footer__control-area__value"
                                        type="text"
                                        id="title"
                                        name="title"
                                        value={socialNetwork.title}
                                        onChange={changeSocialNetwork(socialNetwork.id)}
                                    />
                                    <label htmlFor="link">Link</label>
                                    <input
                                        className="footer__control-area__value"
                                        type="text"
                                        id="link"
                                        name="link"
                                        value={socialNetwork.link}
                                        onChange={changeSocialNetwork(socialNetwork.id)}
                                    />
                                </div>
                        )
                    }
                    <div className="footer__control-area__buttons">
                        <div className="footer__control-area__add-paragraph" onClick={addSocialNetwork}>
                            <img className="footer__control-area__add-paragraph__plus" src={plus} alt="Plus" />
                            <button className="footer__control-area__add-paragraph__btn">
                                <span className="footer__control-area__add-paragraph__btn__text">New Paragraph</span>
                            </button>
                        </div>
                        <div className="footer__control-area__add-paragraph" onClick={saveSocialNetworks}>
                            <button className="footer__control-area__save-paragraph__btn">
                                <span className="footer__control-area__save-paragraph__btn__text">Save</span>
                            </button>
                        </div>
                    </div>
                </div>
                <div className="footer__content">
                    <div className="footer__content__wrapper">
                        {socialNetworks.map((socialNetwork: SocialNetwork, index: number) =>
                            <a
                                key={index}
                                className="footer__content__link"
                                href={socialNetwork.link}
                            >
                                {socialNetwork.title}
                            </a>
                        )}
                    </div>
                </div>
            </div>
        </section>
    );
};

export default Footer;
