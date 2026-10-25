import { APIClient } from '@/api';
import { ImagesNameIdentificators, Layer, LayerImagesIdentificators } from '../layers';

const SECOND_ELEMENT_INDEX = 1;

/**
 * LayersClient is a http implementation of projects API.
 * Exposes all projects-related functionality.
 */
export class LayersClient extends APIClient {
    private readonly ROOT_PATH: string = '/api/v0/collections';
    /** Creates new layer with chosed name and priority */
    public async create(projectId: string, name: string, priority: number): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/layers`;
        const response = await this.http.post(path, JSON.stringify({ name, priority }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Returns list of exist layers for current project*/
    public async list(projectId: string): Promise<Layer[]> {
        const path = `${this.ROOT_PATH}/${projectId}/layers`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const list = await response.json() || [];

        return list.map((layer: Layer) =>
            new Layer(layer.attributes, layer.name, layer.id, layer.priority, layer.projectId, layer.layerImagesCount,
                layer.noneWeight, layer.weights || {}));
    }

    /** Sets the weights of the layer images and of the layer being absent. */
    public async updateWeights(projectId: string, layerId: string, noneWeight: number, weights: Record<string, number>): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/layers/${layerId}/weights`;
        const response = await this.http.put(path, JSON.stringify({ noneWeight, weights }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Updates selected layer priority */
    public async updatePriority(projectId: string, name: string, priority: number): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/layers/${name}`;
        // TODO: string/number is valid json.
        const response = await this.http.patch(path, JSON.stringify({ priority }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Deletes selected layer */
    public async delete(projectId: string, layerId: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/layers/${layerId}`;
        const response = await this.http.delete(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Uploads image on backend */
    public async uploadImage(imageIdentificator: LayerImagesIdentificators): Promise<void> {
        const path = `${this.ROOT_PATH}/${imageIdentificator.projectId}/layers/${imageIdentificator.layerId}/images`;

        const response = await this.http.post(path, JSON.stringify({
            image: imageIdentificator.file,
            name: imageIdentificator.name,
        }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }
    /** Returns list of images names */
    public async getFilesNames(projectId: string, layerId: string): Promise<string[]> {
        const path = `${this.ROOT_PATH}/${projectId}/layers/${layerId}`;

        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const fileNames = await response.json();

        return fileNames || [];
    }

    /** Updates name for selected image */
    public async updateImageName(imageIdentificator: ImagesNameIdentificators): Promise<void> {
        const path = `${this.ROOT_PATH}/${imageIdentificator.projectId}/layers/${imageIdentificator.layerId}/images/${imageIdentificator.name}`;

        const response = await this.http.patch(path, JSON.stringify({
            newName: `${imageIdentificator.newName}.${imageIdentificator.name.split('.')[SECOND_ELEMENT_INDEX]}`,
        }));

        if (!response.ok) {
            await this.handleError(response);
        }
    }

    /** Deletes selected image */
    public async deleteImage(projectId: string, layerId: string, name: string): Promise<void> {
        const path = `${this.ROOT_PATH}/${projectId}/layers/${layerId}/images/${name}`;
        const response = await this.http.delete(path);

        if (!response.ok) {
            await this.handleError(response);
        }
    }
}
