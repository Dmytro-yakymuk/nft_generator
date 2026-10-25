import { User } from '@/accounts';
import { SET_USER } from '@store/actions/accounts';
import { ActionType } from '@/app/store';

export const accountsReducer = (
    accountsState: User = new User(),
    action: ActionType<User> = new ActionType('', new User()),
) => {
    switch (action.type) {
    case SET_USER:
        accountsState = action.payload;
        break;
    default:
        return accountsState;
    };

    return { ...accountsState };
};
