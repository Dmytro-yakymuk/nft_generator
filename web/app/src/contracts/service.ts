import { ContractTerms, ContractsClient } from '@/api/contracts';

export class ContractsService {
    private readonly contracts: ContractsClient;

    constructor(contracts: ContractsClient) {
        this.contracts = contracts;
    }

    /** Returns the data of the transaction that deploys the contract of the project with the given sale terms. */
    public async getContract(projectId: string, terms: ContractTerms = {}): Promise<string> {
        return await this.contracts.getContract(projectId, terms);
    }
}
