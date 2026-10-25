import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { RootState } from '@/app/store';
import { getProject } from '@store/actions/collections';
import { getWhitelistMerkleRoot } from '@/app/store/actions/whitelist';
import { getContractMerkleRoot, setContractMerkleRoot } from '@/ethers/collectionContract';
import { switchNetwork } from '@/ethers/network';
import { metamaskError, notify } from '@/internal/notifications';
import { WhitelistedWallet } from '@/whitelist';

import './index.scss';

const ZERO_ADDRESS = '0x0000000000000000000000000000000000000000';
/** ethers error code: the transaction would revert, here because the wallet is not the contract owner. */
const UNPREDICTABLE_GAS_LIMIT = 'UNPREDICTABLE_GAS_LIMIT';

type ContractMerkleRootProps = {
    projectId: string;
    /** The shown whitelist page: the roots are compared again whenever it is reloaded after a change. */
    wallets: WhitelistedWallet[];
};

/** Shows whether the deployed contract holds the Merkle root of the current whitelist and lets the owner update it:
 * after the whitelist changes, the wallets can mint through the whitelist only with the new root in the contract. */
export const ContractMerkleRoot: React.FC<ContractMerkleRootProps> = ({ projectId, wallets }) => {
    const dispatch = useDispatch();
    const project = useSelector((state: RootState) => state.collectionsReducer.project);
    const [whitelistRoot, setWhitelistRoot] = useState<string>('');
    const [contractRoot, setContractRoot] = useState<string>('');
    const [isUpdating, setIsUpdating] = useState<boolean>(false);
    const contractAddress = project.id === projectId ? project.contractAddress : ZERO_ADDRESS;
    const isDeployed = contractAddress.toLowerCase() !== ZERO_ADDRESS;

    useEffect(() => {
        dispatch(getProject(projectId));
    }, [projectId]);

    useEffect(() => {
        if (!isDeployed) {
            return;
        }
        (async() => {
            try {
                const [root, onChainRoot] = await Promise.all([
                    getWhitelistMerkleRoot(projectId),
                    getContractMerkleRoot(contractAddress),
                ]);
                setWhitelistRoot(root.toLowerCase());
                setContractRoot(onChainRoot.toLowerCase());
            } catch (error) {
                notify('Could not compare the whitelist with the contract. Try later.');
            }
        })();
    }, [isDeployed, contractAddress, wallets]);

    /** Sends setMerkleRoot from the owner wallet in MetaMask. */
    const updateContractRoot = async() => {
        setIsUpdating(true);
        try {
            await switchNetwork();
            await setContractMerkleRoot(contractAddress, whitelistRoot);
            // a public node may still return the previous state for a few seconds, the receipt is enough.
            setContractRoot(whitelistRoot);
            notify('The contract now holds the current whitelist.', 'success');
        } catch (error: any) {
            if (error.code === UNPREDICTABLE_GAS_LIMIT) {
                notify('Only the wallet that deployed the contract can update its whitelist.');
            } else {
                metamaskError(error);
            }
        } finally {
            setIsUpdating(false);
        }
    };

    if (!isDeployed || !whitelistRoot || !contractRoot) {
        return null;
    }

    if (whitelistRoot === contractRoot) {
        return (
            <div className="contract-merkle-root contract-merkle-root--current">
                The contract holds the current whitelist.
            </div>
        );
    }

    return (
        <div className="contract-merkle-root contract-merkle-root--outdated">
            <span className="contract-merkle-root__text">
                The whitelist has changed since the contract got it: the listed wallets cannot mint through the whitelist
                until the contract gets the new Merkle root.
            </span>
            <button
                className="contract-merkle-root__btn"
                disabled={isUpdating}
                onClick={updateContractRoot}
            >
                {isUpdating ? 'Updating...' : 'Update contract'}
            </button>
        </div>
    );
};
