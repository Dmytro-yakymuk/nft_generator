import { APIClient } from '@/api';

/** Sale terms the owner sets for the collection contract; empty fields take the server defaults. */
export interface ContractTerms {
    name?: string;
    /** Price of one token, ETH. */
    price?: string;
    /** Maximum number of tokens in one transaction. */
    maxPerMint?: number;
}

/** Client of the contracts API: the data of the transaction that deploys a collection contract. */
export class ContractsClient extends APIClient {
    private readonly ROOT_PATH: string = '/api/v0/contracts';

    /** Returns the data of the transaction that deploys the contract of the project with the given sale terms. */
    public async getContract(projectId: string, terms: ContractTerms = {}): Promise<string> {
        const query = new URLSearchParams();
        terms.name && query.set('name', terms.name);
        terms.price && query.set('price', terms.price);
        terms.maxPerMint && query.set('maxPerMint', String(terms.maxPerMint));
        const path = `${this.ROOT_PATH}/data/${projectId}${query.toString() ? `?${query}` : ''}`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        return await response.json();
    }
}
