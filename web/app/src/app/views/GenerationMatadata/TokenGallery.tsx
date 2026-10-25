import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { RootState } from '@/app/store';
import { getLayers } from '@/app/store/actions/layers';
import { Layer } from '@/layers';
import { NONE_NAME, percentOf, rarityTier } from '@/layers/rarity';
import '@/layers/rarity.scss';

const TOKENS_SHOWN = 12;
const FIRST_TOKEN = 0;
const PERCENT_DIGITS = 1;

/** Attribute of the token metadata: the layer and its variant. */
interface Attribute {
    // eslint-disable-next-line camelcase
    trait_type: string;
    value: string;
}

/** Token of the generated collection as the contract's base URI serves it. */
interface Token {
    id: number;
    name: string;
    image: string;
    attributes: Attribute[];
}

/** Returns the share in percent of the images that get the variant, by the saved weights of the layer. */
const variantShare = (layers: Layer[], attribute: Attribute): number => {
    const layer = layers.find((candidate: Layer) => candidate.name === attribute.trait_type);
    if (!layer) {
        return NaN;
    }
    const weights = Object.keys(layer.weights).map((image: string) => layer.weights[image]);
    const total = weights.reduce((sum: number, weight: number) => sum + weight, layer.noneWeight);
    if (attribute.value === NONE_NAME) {
        return percentOf(layer.noneWeight, total);
    }
    const image = Object.keys(layer.weights).find((name: string) => name.replace(/\.[^.]+$/, '') === attribute.value);

    return image ? percentOf(layer.weights[image], total) : NaN;
};

/** Returns the attributes of the token one per line. */
const describe = (token: Token): string =>
    token.attributes.map((attribute: Attribute) => `${attribute.trait_type}: ${attribute.value}`).join('\n');

/** Returns the rarest trait of the token by the saved weights, if the layers are known. */
const rarestTrait = (layers: Layer[], token: Token) => token.attributes
    .map((attribute: Attribute) => ({ attribute, share: variantShare(layers, attribute) }))
    .filter(({ share }) => !isNaN(share))
    .sort((a, b) => a.share - b.share)[FIRST_TOKEN];

/** Loads the metadata of the token from the address the contract's base URI points to. */
const loadToken = async(projectId: string, id: number): Promise<Token | null> => {
    const response = await fetch(`${window.location.origin}/metadata/${projectId}/${id}`);

    return response.ok ? { id, ...await response.json() } : null;
};

const isToken = (token: Token | null): token is Token => token !== null;

/** The first tokens of the generated collection with their rarest trait. */
export const TokenGallery: React.FC<{ projectId: string; count: number }> = ({ projectId, count }) => {
    const dispatch = useDispatch();
    const layers = useSelector((state: RootState) => state.layersReducer.layers);
    const [tokens, setTokens] = useState<Token[]>([]);

    useEffect(() => {
        dispatch(getLayers(projectId));
    }, [projectId]);

    useEffect(() => {
        const ids = Array.from({ length: Math.min(TOKENS_SHOWN, count) }, (_, i) => FIRST_TOKEN + i);
        Promise.all(ids.map((id: number) => loadToken(projectId, id)))
            .then((loaded) => setTokens(loaded.filter(isToken)));
    }, [projectId, count]);

    if (!tokens.length) {
        return null;
    }

    return (
        <div className="token-gallery">
            <h2 className="token-gallery__title">First tokens</h2>
            <div className="token-gallery__grid">
                {tokens.map((token: Token) => {
                    const rarest = rarestTrait(layers, token);
                    const tier = rarest && rarityTier(rarest.share);

                    return (
                        <figure className="token-gallery__token" key={token.id} title={describe(token)}>
                            <img src={token.image} alt={token.name} />
                            <figcaption>
                                <span className="token-gallery__name">#{token.id}</span>
                                {rarest &&
                                    <span className="token-gallery__trait">
                                        <span className={`rarity rarity--${tier.toLowerCase()}`}>{tier}</span>
                                        {rarest.attribute.value} {rarest.attribute.trait_type}, {rarest.share.toFixed(PERCENT_DIGITS)}%
                                    </span>
                                }
                            </figcaption>
                        </figure>
                    );
                })}
            </div>
        </div>
    );
};
