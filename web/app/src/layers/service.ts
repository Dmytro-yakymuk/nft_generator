import { LayersClient } from '@/api/layers';
import { ImagesNameIdentificators, Layer, LayerIdentificators, LayerImagesIdentificators } from '.';

export class LayersService {
    private readonly layers: LayersClient;

    constructor(layers: LayersClient) {
        this.layers = layers;
    }

    /** Creates new layer with chosed name and priority */
    public async create(projectId: string, name: string, priority: number): Promise<void> {
        await this.layers.create(projectId, name, priority);
    }

    /** Returns list of exist layers for current project*/
    public async list(projectId: string): Promise<Layer[]> {
        return await this.layers.list(projectId);
    }

    /** Sets the weights of the layer images and of the layer being absent. */
    public async updateWeights(projectId: string, layerId: string, noneWeight: number, weights: Record<string, number>): Promise<void> {
        await this.layers.updateWeights(projectId, layerId, noneWeight, weights);
    }

    /** Updates selected layer priority */
    public async updatePriority(projectId: string, name: string, priority: number): Promise<void> {
        await this.layers.updatePriority(projectId, name, priority);
    }

    /** Deletes selected layer */
    public async delete(projectId: string, layerId: string): Promise<void> {
        await this.layers.delete(projectId, layerId);
    }

    /** Uploads image on backend */
    public async uploadImage(imageIdentificator: LayerImagesIdentificators): Promise<void> {
        await this.layers.uploadImage(imageIdentificator);
    }

    /** Returns  list io images names */
    public async getFilesNames(projectId: string, layerId: string): Promise<string[]> {
        return await this.layers.getFilesNames(projectId, layerId);
    }

    /** Updates name for selected image */
    public async updateImageName(imageIdentificator: ImagesNameIdentificators): Promise<void> {
        await this.layers.updateImageName(imageIdentificator);
    }

    /** Deletes selected image */
    public async deleteImage(projectId: string, layerId: string, name: string): Promise<void> {
        await this.layers.deleteImage(projectId, layerId, name);
    }
}
