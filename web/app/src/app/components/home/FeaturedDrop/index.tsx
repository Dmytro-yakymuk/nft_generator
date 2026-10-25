import { Fragment } from 'react';
import Slider from 'react-slick';

import image2 from '@static/images/featured-drop/2.webp';
import image3 from '@static/images/featured-drop/3.webp';
import image4 from '@static/images/featured-drop/4.webp';

import './index.scss';

export const FeaturedDrop = () => {
    /** Images list for slider. */
    const images: string[] = [image2, image3, image4];
    const ARRAY_LENGTH = 10;

    /** Settings for images slider. */
    const settings = {
        adaptiveHeight: true,
        infinite: true,
        dots: false,
        speed: 500,
        arrows: false,
        autoplay: true,
        slidesToShow: 1,
        slidesToScroll: 1,
    };
    const words = new Array(ARRAY_LENGTH).fill('Featured Drop');

    return (
        <section className="featured-drop" id="featured-drop">
            <div className="wrapper">
                <div className="featured-drop__title">
                    {
                        words.map(
                            (word: string, index: number) =>
                                <Fragment key={index}>
                                    <span className="white-title">{word}</span>
                                    <span className="black-title">{word}</span>
                                </Fragment>
                        )
                    }
                </div>
                <Slider {...settings} className="slider">
                    {images.map((image: string, index: number) =>
                        <img
                            src={image}
                            key={index}
                            alt="image"
                            loading="lazy"
                            className="featured-drop__image"
                        />
                    )}
                </Slider>
            </div>
        </section>
    );
};
