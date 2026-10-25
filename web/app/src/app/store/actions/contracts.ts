import { Dispatch } from 'redux';

import { ContractTerms, ContractsClient } from '@/api/contracts';
import { ContractsService } from '@/contracts/service';

/** SET_C0NTRACT is action type for direct manipulation setContract action of store. */
export const SET_C0NTRACT = 'SET_C0NTRACT';

export const setContract = (contractSignature: string) => ({
    type: SET_C0NTRACT,
    payload: contractSignature,
});

const contractsClient = new ContractsClient();
const contractsService = new ContractsService(contractsClient);

/** Requests the data of the transaction that deploys the project contract with the given sale terms and sets it into state. */
export const getContract = (projectId: string, terms: ContractTerms = {}) => async(dispatch: Dispatch) => {
    const contractSignature = await contractsService.getContract(projectId, terms);
    dispatch(setContract(contractSignature));
};

