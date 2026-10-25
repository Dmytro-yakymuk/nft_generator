const DEFAULT_VALUE = 0;

/** Exposes layer domain entity. */
export class Layer {
    /** At least, contains name, id and files list. */
    public constructor(
        public attributes: any = null,
        public name: string = 'Create layer',
        public id: string = '',
        public priority: number = DEFAULT_VALUE,
        public projectId: string = '',
        public layerImagesCount: number = DEFAULT_VALUE,
        /** Weight of the layer being absent from an image; zero makes the layer mandatory. */
        public noneWeight: number = DEFAULT_VALUE,
        /** Weights of the layer images by file name. */
        public weights: Record<string, number> = {},
    ) { };
};

/** Exposes fields for layers manipulation */
export class LayerIdentificators {
    public constructor(
        public projectId: string,
        public name: string,
        public priority: number
    ) { }
}

/** Exposes fields for images manipulation */
export class LayerImagesIdentificators {
    public constructor(
        public projectId: string = '',
        public layerId: string = '',
        public name: string = '',
        public file: string = ''
    ) { }
}
/** Exposes fields for images names manipulation */
export class ImagesNameIdentificators {
    public constructor(
        public projectId: string = '',
        public layerId: string = '',
        public name: string = '',
        public newName: string = ''
    ) { }
}
