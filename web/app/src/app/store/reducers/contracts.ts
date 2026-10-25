import { ActionType } from '@/app/store';
import { SET_C0NTRACT } from '../actions/contracts';

/** Exposes contract state. */
class ContractsState {
    constructor(
        public contractSignature: string = ''
    ) { }
}

export const contractsReducer = (
    contractsState: ContractsState = new ContractsState(),
    action: ActionType<string> = new ActionType('', ''),
) => {
    switch (action.type) {
    case SET_C0NTRACT:
        contractsState.contractSignature = action.payload;
        break;
    default:
        break;
    }

    return { ...contractsState };
};
