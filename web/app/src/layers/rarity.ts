import { Layer } from '.';

/** Weight of an image whose weight was never set; equal weights mean equal chances (as on the server). */
export const DEFAULT_WEIGHT = 100;
/** Largest weight the server accepts. */
export const MAX_WEIGHT = 1000000;
/** Attribute value of a layer that is absent from an image. */
export const NONE_NAME = 'None';

const PERCENT = 100;
const NONE = 0;
const LAST = 1;

/** Rarity tier of a variant by the share of images that get it: the tier holds shares up to maxPercent. */
export const RARITY_TIERS = [
    { name: 'Legendary', maxPercent: 1 },
    { name: 'Epic', maxPercent: 5 },
    { name: 'Rare', maxPercent: 10 },
    { name: 'Uncommon', maxPercent: 20 },
    { name: 'Common', maxPercent: PERCENT },
];

/** Returns the rarity tier for the share of images in percent. */
export const rarityTier = (percent: number): string =>
    (RARITY_TIERS.find((tier) => percent <= tier.maxPercent) || RARITY_TIERS[RARITY_TIERS.length - LAST]).name;

/** Variant of a layer: an image or, for an optional layer, its absence (image is null). */
export class Variant {
    public constructor(
        public name: string,
        public image: string | null,
        public weight: number,
    ) { }
}

/** Returns the variants of the layer in the server order: images by file name, then the absence. */
export const layerVariants = (layer: Layer): Variant[] => {
    const variants = Object.keys(layer.weights).sort()
        .map((image) => new Variant(image.replace(/\.[^.]+$/, ''), image, layer.weights[image]));
    if (layer.noneWeight > NONE) {
        variants.push(new Variant(NONE_NAME, null, layer.noneWeight));
    }

    return variants;
};

/** Returns the share of the weight in percent of the total. */
export const percentOf = (weight: number, total: number): number => total > NONE ? weight / total * PERCENT : NONE;

/** Draws a variant with probability proportional to its weight. */
export const drawVariant = (variants: Variant[], random: () => number = Math.random): Variant => {
    const total = variants.reduce((sum, variant) => sum + variant.weight, NONE);
    let point = random() * total;
    for (const variant of variants) {
        point -= variant.weight;
        if (point < NONE) {
            return variant;
        }
    }

    return variants[variants.length - LAST];
};
