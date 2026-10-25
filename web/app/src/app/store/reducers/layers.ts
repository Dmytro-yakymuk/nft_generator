import { REFRESH_LAYERS, SET_ACTIVE_LAYER, SET_FILE_NAMES, SET_LAYER, SET_PRIORITY } from '@/app/store/actions/layers';
import { Layer } from '@/layers';
import { ActionType } from '@/app/store';

const SECOND_ELEMENT_INDEX = 1;

/** Exposes layers entity state. */
class LayersState {
    constructor(
        public layers: Layer[] = [],
        public activeLayer: Layer = new Layer(),
        public fileNames: string[] = [],
    ) { };
};

export const layersReducer = (
    layersState: LayersState = new LayersState(),
    action: ActionType<any> = new ActionType('', ''),
) => {
    switch (action.type) {
    case SET_LAYER:
        layersState.layers = [...action.payload];
        layersState.activeLayer = layersState.layers[layersState.layers.length - SECOND_ELEMENT_INDEX] || new Layer();
        layersState.fileNames = [];
        break;
    case SET_PRIORITY:
        layersState.layers = [...action.payload];
        break;
    case REFRESH_LAYERS:
        layersState.layers = [...action.payload];
        layersState.activeLayer = layersState.layers.find((layer: Layer) => layer.id === layersState.activeLayer.id)
            || layersState.layers[layersState.layers.length - SECOND_ELEMENT_INDEX] || new Layer();
        break;
    case SET_ACTIVE_LAYER:
        layersState.activeLayer = { ...layersState.layers[action.payload] };
        layersState.fileNames = [];
        break;
    case SET_FILE_NAMES:
        layersState.fileNames = [...action.payload];
        break;
    default:
        return layersState;
    }

    return { ...layersState };
};
