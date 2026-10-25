import { useEffect, useState } from 'react';

import { getFirstTokenURI, setContractBaseURI } from '@/ethers/collectionContract';
import { switchNetwork } from '@/ethers/network';
import { metamaskError, notify } from '@/internal/notifications';

const ZERO_ADDRESS = '0x0000000000000000000000000000000000000000';
/** ethers error code: the transaction would revert, here because the wallet is not the contract owner. */
const UNPREDICTABLE_GAS_LIMIT = 'UNPREDICTABLE_GAS_LIMIT';

type ContractBaseURIProps = {
    contractAddress: string;
    metadataCid: string;
    gateway: string;
};

/** Shows where the deployed contract takes token metadata from and lets the owner point it to the collection published
 * to IPFS: a contract deployed before the publication still returns the platform address in tokenURI. */
export const ContractBaseURI: React.FC<ContractBaseURIProps> = ({ contractAddress, metadataCid, gateway }) => {
    const baseURI = `ipfs://${metadataCid}/`;
    const [tokenURI, setTokenURI] = useState<string | null>(null);
    const [isUpdating, setIsUpdating] = useState<boolean>(false);
    const [isUpdated, setIsUpdated] = useState<boolean>(false);
    const isDeployed = contractAddress.toLowerCase() !== ZERO_ADDRESS;

    useEffect(() => {
        if (!isDeployed) {
            return;
        }
        getFirstTokenURI(contractAddress)
            .then(setTokenURI)
            .catch(() => notify('Could not read the contract. Try later.'));
    }, [contractAddress, isUpdated]);

    /** Sends setBaseURI from the owner wallet in MetaMask. */
    const pointToIPFS = async() => {
        setIsUpdating(true);
        try {
            await switchNetwork();
            await setContractBaseURI(contractAddress, baseURI);
            setIsUpdated(true);
            notify('The contract now takes token metadata from IPFS.', 'success');
        } catch (error: any) {
            if (error.code === UNPREDICTABLE_GAS_LIMIT) {
                notify('Only the wallet that deployed the contract can change its base URI.');
            } else {
                metamaskError(error);
            }
        } finally {
            setIsUpdating(false);
        }
    };

    if (!isDeployed || tokenURI === null) {
        return null;
    }

    // without a minted token the success of the transaction is all that can be shown.
    const isOnIPFS = tokenURI.startsWith(baseURI) || isUpdated && !tokenURI;
    if (isOnIPFS) {
        return (
            <div className="contract-base-uri contract-base-uri--ipfs">
                <span>
                    The contract takes token metadata from IPFS
                    {tokenURI && <>: tokenURI(0) = <a href={`${gateway}/ipfs/${tokenURI.slice('ipfs://'.length)}`} target="_blank" rel="noreferrer">{tokenURI}</a></>}.
                </span>
            </div>
        );
    }

    return (
        <div className="contract-base-uri contract-base-uri--outdated">
            <span>
                {tokenURI
                    ? `The contract still returns tokenURI(0) = ${tokenURI}.`
                    : 'No token is minted yet, so the contract base URI cannot be read.'}
            </span>
            <button className="contract-base-uri__btn" disabled={isUpdating} onClick={pointToIPFS}>
                {isUpdating ? 'Updating...' : 'Point the contract to IPFS'}
            </button>
        </div>
    );
};
