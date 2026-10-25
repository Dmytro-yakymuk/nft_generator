import { Dispatch } from 'redux';

import { User, UserEditFields } from '@/accounts';
import { UsersClient } from '@/api/users';
import { UsersService } from '@/accounts/service';
import { session } from '@/accounts/session';

export const SET_USER = 'SET_USER';

/** An action setUser contains type and payload data for setting user. */
export const setUser = (user: User) => ({
    type: SET_USER,
    payload: user,
});

const usersClient = new UsersClient();
const usersService = new UsersService(usersClient);

/** Thunk middleware that gets user from api and sets it into state */
export const getUser = () => async function(dispatch: Dispatch) {
    const user = await usersService.getUser();
    dispatch(setUser(user));
};

/** Thunk middleware that updates user from api and sets updated user into state */
export const updateUser = (editFields: UserEditFields) => async function(dispatch: Dispatch) {
    await usersService.updateuser(editFields);
    const user = await usersService.getUser();
    dispatch(setUser(user));
};

/** Thunk middleware that logs out and clears the user in state */
export const logout = () => async function(dispatch: Dispatch) {
    await usersService.logout();
    session.end();
    dispatch(setUser(new User()));
};
