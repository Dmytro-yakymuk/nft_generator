import { useMemo } from 'react';
import MetaMaskOnboarding from '@metamask/onboarding';
import { useHistory, useLocation } from 'react-router-dom';

import { NotFoundError } from '@/api';
import { UsersClient } from '@/api/users';
import { UsersService } from '@/accounts/service';
import { session } from '@/accounts/session';
import { RoutesConfig } from '@/app/routes';
import { SignedMessage } from '@/ethers';
import { notify } from '@/internal/notifications';

/**
 * Hook useMetamaskLogin returns the login function: a new wallet is registered first, then signs the nonce.
 * From the home page the user goes to the collections, elsewhere stays on the page.
 */
export const useMetamaskLogin = () => {
    const ROOT_PATH = '/';
    const history = useHistory();
    const location = useLocation();
    const onboarding = useMemo(() => new MetaMaskOnboarding(), []);
    const service = useMemo(() => new UsersService(new UsersClient()), []);
    /** The ethers.js service is loaded only when the user logs in, so the first page does not wait for it. */
    const loadEthersService = async() => (await import('@/plugins/service')).ServicePlugin.create();

    /** Signs the nonce of the wallet and logs in. */
    const signNonce = async() => {
        const ethersService = await loadEthersService();
        const address = await ethersService.getWallet();
        const nonce = await service.getNonce(address);
        const signature = await ethersService.signMessage(nonce);
        await service.login(new SignedMessage(nonce, signature));
        session.start();
        if (location.pathname === ROOT_PATH) {
            history.push(RoutesConfig.Collections.path);
        }
    };

    return async() => {
        if (!MetaMaskOnboarding.isMetaMaskInstalled()) {
            onboarding.startOnboarding();

            return;
        }

        await window.ethereum.request({
            method: 'eth_requestAccounts',
        });
        try {
            await signNonce();
        } catch (error: any) {
            if (!(error instanceof NotFoundError)) {
                notify('Something went wrong');

                return;
            }
            try {
                const ethersService = await loadEthersService();
                const signature = await ethersService.signMessage('Register with metamask');
                await service.register(signature);
                await signNonce();
            } catch (error: any) {
                notify('Something went wrong');
            }
        }
    };
};
