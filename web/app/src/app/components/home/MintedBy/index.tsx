import mintedBy from '@static/images/minted.webp';

import './index.scss';

/** The image is the largest element of the first screen, so the browser fetches it before the rest. */
const HIGH_PRIORITY = { fetchpriority: 'high' };

export const MintedBy = () =>
    <section className="minted-by">
        <img className="minted-by__image" src={mintedBy} alt="Minted by image" {...HIGH_PRIORITY} />
        <span className="minted-by__title">minted by generator.</span>
    </section>;

