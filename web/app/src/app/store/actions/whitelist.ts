import { Dispatch } from 'redux';

import { WhitelistClient } from '@/api/whitelist';
import { WhitelistService } from '@/whitelist/service';
import { WhitelistPage, WhitelistedWallet } from '@/whitelist';

export const SET_ADDRESS_TO_WHITELIST = 'SET_ADDRESS_TO_WHITELIST';
export const DELETE_FROM_WHITELIST = 'DELETE_FROM_WHITELIST';
export const SET_WHITELIST = 'SET_WHITELIST';

/** An action setWhitelist contains type and payload data for sets whitelist. */
export const setWhitelist = (whitelist: WhitelistPage) => ({
    type: SET_WHITELIST,
    payload: whitelist,
});

/** An action setToWhitelist contains type and payload data for sets wallet address to whitelist. */
export const setToWhitelist = (whitelistedWallet: WhitelistedWallet) => ({
    type: SET_ADDRESS_TO_WHITELIST,
    payload: whitelistedWallet,
});

/** An action deleteFromWhitelist contains type and payload data for deletes wallet address from whitelist. */
export const deleteFromWhitelist = (projectId: string, walletAddress: string) => ({
    type: DELETE_FROM_WHITELIST,
    payload: { projectId, walletAddress },
});

const whitelistClient = new WhitelistClient();
const whitelistService = new WhitelistService(whitelistClient);

/** Thunk middleware that requests all projects and sets whitelist to store. */
export const getWhitelist = (projectId: string, currentPage: number, name: string) => async function(dispatch: Dispatch) {
    const whitelist = await whitelistService.getWhitelist(projectId, currentPage, name);
    dispatch(setWhitelist(whitelist));
};

/** Thunk middleware that adds wallet address and sets updated whitelist list to store. */
export const addWalletToWhitelist = (projectId: string, name: string, walletAddress: string) => async function(dispatch: Dispatch) {
    await whitelistService.addWalletToWhitelist(projectId, name, walletAddress);
    const whitelist = await whitelistService.getWhitelist(projectId);
    dispatch(setWhitelist(whitelist));
};

/** Thunk middleware that deletes wallet address and sets updated whitelist list to store. */
export const deleteWalletFromWhitelist = (projectId: string, walletAddress: string, currentPage: number) => async function(dispatch: Dispatch) {
    await whitelistService.deleteWalletFromWhitelist(projectId, walletAddress);
    const whitelist = await whitelistService.getWhitelist(projectId, currentPage);
    dispatch(setWhitelist(whitelist));
};

/** Thunk middleware that uploads whitelist excel file and sets updated whitelist list to store. */
export const uploadWhitelistFile = (projectId: string, file: string) => async function(dispatch: Dispatch) {
    await whitelistService.uploadWhitelistFile(projectId, file);
    const whitelist = await whitelistService.getWhitelist(projectId);
    dispatch(setWhitelist(whitelist));
};

/** Returns the Merkle root of the whitelist that the deployed contract must hold. */
export const getWhitelistMerkleRoot = async(projectId: string): Promise<string> => await whitelistService.getMerkleRoot(projectId);
