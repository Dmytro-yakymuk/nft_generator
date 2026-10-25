import { Link } from 'react-router-dom';

import { RoutesConfig } from '@/app/routes';

import chart from '@static/images/howItWorks/chart.png';
import crown from '@static/images/howItWorks/crown.png';
import diamond from '@static/images/howItWorks/diamond.png';
import file from '@static/images/howItWorks/file.png';
import symbolChart from '@static/images/howItWorks/symbol_chart.png';
import symbolDiamond from '@static/images/howItWorks/symbol_diamond.png';
import symbolWorld from '@static/images/howItWorks/symbol_world.png';
import world from '@static/images/howItWorks/world.png';

import './index.scss';

/** Exposes a step of launching a collection. */
class Step {
    /** Contains label, description and image. */
    constructor(
        public label: string,
        public description: string,
        public image: string,
    ) { };
};

/** Exposes feature entity of project. */
class Feature {
    /** At least, contains label and image. */
    constructor(
        public label: string,
        public image: string,
    ) { };
};

export const HowItWorks: React.FC = () => {
    /** Features list. */
    const features: Feature[] = [
        new Feature('Contracts', file),
        new Feature('Websites', world),
        new Feature('Collections', diamond),
        new Feature('Analytics', chart),
        new Feature('IPFS', crown),
    ];

    /** Steps from image layers to a published mint site. */
    const steps: Step[] = [
        new Step('Generating a Collection', 'Upload image layers, set their order and generate unique images with ERC-721 metadata.', symbolDiamond),
        new Step('Publishing to IPFS', 'Publish the images and metadata to IPFS: token URIs depend on the content, not on our server.', symbolWorld),
        new Step('Contract Creation', 'Deploy the ERC-721 contract from MetaMask and keep the whitelist as a Merkle root.', symbolChart),
        new Step('Creating a Website', 'Build a mint site from sections, publish it for buyers and follow the sales on the dashboard.', symbolWorld),
    ];

    return <section className="works">
        <div className="works__wrapper">

            <h1 className="works__title" id="features">
                Features
            </h1>
            <ul className="works__features">
                {features.map((feature: Feature, index: number) =>
                    <li
                        className="works__features__item"
                        key={index}
                    >
                        <img
                            className="works__features__item__image"
                            src={feature.image}
                            alt={`${feature.label}`}
                        />
                        <span className="works__features__item__label">
                            {feature.label}
                        </span>
                    </li>
                )}
            </ul>
            <h1 className="works__title" id="how-it-works">
                How It Works
            </h1>
            <div className="works__main">
                {steps.map((step: Step, index: number) =>
                    <div className="works__main__item" key={index}>
                        <div className="works__main__item__logo">
                            <img src={step.image} alt={step.label} />
                        </div>
                        <div className="works__main__item__divider" />
                        <span className="works__main__item__label">
                            {step.label}
                        </span>
                        <p className="works__main__item__description">
                            {step.description}
                        </p>
                    </div>
                )}
            </div>
            <Link className="works__try-it" to={RoutesConfig.Collections.path}>
                Try It Out
            </Link>
        </div>
    </section>;
};
