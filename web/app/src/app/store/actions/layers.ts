import { LayersClient } from '@/api/layers';
import { ImagesNameIdentificators, Layer, LayerIdentificators, LayerImagesIdentificators } from '@/layers';
import { LayersService } from '@/layers/service';
import { Dispatch } from 'redux';

/** SET_LAYER is action type for direct manipulation setLayer action of store. */
export const SET_LAYER = 'SET_LAYER';
/** UPDATE_LAYERS_FILES_COUNT is action type for direct manipulation updateLayersFIlesCount action of store. */
export const UPDATE_LAYERS_FILES_COUNT = 'UPDATE_LAYERS_FILES_COUNT';
/** DELETE_LAYER is action type for direct manipulation deleteLayer action of store. */
export const DELETE_LAYER = 'DELETE_LAYER';
/** SET_LAYER is action type for direct manipulation setFiles action of store. */
export const ADD_FILES = 'ADD_FILES';
/** DELETE_FILE is action type for direct manipulation deleteFile action of store. */
export const DELETE_FILE = 'DELETE_FILE';
/** CHANGE_FILE_NAME is action type for direct manipulation changeFileName action of store.*/
export const CHANGE_FILE_NAME = 'CHANGE_FILE_NAME';
/** SET_ACTIVE_LAYER is action type for direct manipulation setActivelayer action of store.*/
export const SET_ACTIVE_LAYER = 'SET_ACTIVE_LAYER';
/** EXCHANGE_LAYERS is action type for direct manipulation exchangeLayers action of store.*/
export const EXCHANGE_LAYERS = 'EXCHANGE_LAYERS';
/** SET_FILE_NAMES is action type for direct manipulation setFileNames action of store.*/
export const SET_FILE_NAMES = 'SET_FILE_NAMES';
/** UPDATE_PRIORITY is action type for direct manipulation setFileNames action of store.*/
export const SET_PRIORITY = 'SET_PRIORITY';
/** REFRESH_LAYERS replaces the layers but keeps the active layer selected. */
export const REFRESH_LAYERS = 'REFRESH_LAYERS';

export const setLayers = (layers: Layer[]) => ({
    type: SET_LAYER,
    payload: layers,
});
export const updateLayersFilesCount = (layers: Layer[]) => ({
    type: REFRESH_LAYERS,
    payload: layers,
});
export const setPriority = (layers: Layer[]) => ({
    type: SET_PRIORITY,
    payload: layers,
});
export const addFiles = (files: File[]) => ({
    type: ADD_FILES,
    payload: files,
});
export const setFileNames = (names: string[]) => ({
    type: SET_FILE_NAMES,
    payload: names,
});
export const deleteFile = (index: number) => ({
    type: DELETE_FILE,
    payload: index,
});
export const setActivelayer = (index: number) => ({
    type: SET_ACTIVE_LAYER,
    payload: index,
});
export const changeFileName = (index: number, name: string) => ({
    type: CHANGE_FILE_NAME,
    payload: { index, name },
});

const layersClient = new LayersClient();
const layersService = new LayersService(layersClient);

/** Returns list of exist layers for current project*/
export const getLayers = (projectId: string) => async(dispatch: Dispatch) => {
    const layers = await layersService.list(projectId);
    dispatch(setLayers(layers));
};

/** Creates new layer with chosed name and priority */
export const addLayer = (projectId: string, name: string, priority: number) => async(dispatch: Dispatch) => {
    await layersService.create(projectId, name, priority);
    const layers = await layersService.list(projectId);

    dispatch(setLayers(layers));
};

/** Sets the weights of the layer images and of the layer being absent. */
export const updateWeights = (projectId: string, layerId: string, noneWeight: number, weights: Record<string, number>) =>
    async(dispatch: Dispatch) => {
        await layersService.updateWeights(projectId, layerId, noneWeight, weights);
        const layers = await layersService.list(projectId);
        dispatch(updateLayersFilesCount(layers));
    };

/** Updates selected layer priority */
export const updatePriority = (projectId: string, name: string, priority: number) => async(dispatch: Dispatch) => {
    await layersService.updatePriority(projectId, name, priority);
    const layers = await layersService.list(projectId);
    dispatch(setPriority(layers));
};

/** Deletes selected layer */
export const deleteLayer = (layerId: string, projectId: string) => async(dispatch: Dispatch) => {
    await layersService.delete(projectId, layerId);
    const layers = await layersService.list(projectId);
    dispatch(setLayers(layers));
};

/** Uploads image on backend */
export const uploadImage = (imageIdentificator: LayerImagesIdentificators) => async(dispatch: Dispatch) => {
    await layersService.uploadImage(imageIdentificator);
    const names = await layersService.getFilesNames(imageIdentificator.projectId, imageIdentificator.layerId);
    const layers = await layersService.list(imageIdentificator.projectId);
    dispatch(setFileNames(names));
    dispatch(updateLayersFilesCount(layers));
};

/** Returns  list io images names */
export const getFilesNames = (projectId: string, layerId: string) => async(dispatch: Dispatch) => {
    const names = await layersService.getFilesNames(projectId, layerId);
    dispatch(setFileNames(names));
};

/** Updates name for selected image */
export const updateImageName = (imageIdentificator: ImagesNameIdentificators) => async(dispatch: Dispatch) => {
    await layersService.updateImageName(imageIdentificator);
    const names = await layersService.getFilesNames(imageIdentificator.projectId, imageIdentificator.layerId);
    const layers = await layersService.list(imageIdentificator.projectId);
    dispatch(setFileNames(names));
    dispatch(updateLayersFilesCount(layers));
};

/** Deletes selected image */
export const deleteImage = (projectId: string, layerId: string, name: string) => async(dispatch: Dispatch) => {
    await layersService.deleteImage(projectId, layerId, name);
    const names = await layersService.getFilesNames(projectId, layerId);
    const layers = await layersService.list(projectId);
    dispatch(setFileNames(names));
    dispatch(updateLayersFilesCount(layers));
};
