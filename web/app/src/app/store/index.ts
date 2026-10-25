import thunk from 'redux-thunk';
import { applyMiddleware, combineReducers, createStore } from 'redux';
import { dashboardReducer } from '@/app/store/reducers/dashboard';
import { layersReducer } from '@/app/store/reducers/layers';
import { collectionsReducer } from '@/app/store/reducers/collection';
import { whitelistReducer } from '@/app/store/reducers/whitelist';
import { accountsReducer } from '@/app/store/reducers/accounts';
import { contractsReducer } from '@/app/store/reducers/contracts';

/** Common generic action type */
export class ActionType<PayloadType> {
    constructor(
        public type: string,
        public payload: PayloadType
    ) { }
}

const reducer = combineReducers({
    dashboardReducer,
    layersReducer,
    collectionsReducer,
    whitelistReducer,
    accountsReducer,
    contractsReducer,
});

export const store = createStore(reducer, applyMiddleware(thunk));

export type RootState = ReturnType<typeof store.getState>;
