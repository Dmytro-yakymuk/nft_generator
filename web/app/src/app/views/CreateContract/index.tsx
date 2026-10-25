import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { NavLink, useParams } from 'react-router-dom';

import { Navbar } from '@components/common/Navbar';
import { ProjectNavigation } from '@/app/components/common/ProjectNavigation';
import StyledSelect from '@/app/components/CreateContract/StyledSelect';
import { Popup } from '@components/common/Popup';

import pinkLogo from '@static/images/logo.png';

import { getContract, setContract } from '@store/actions/contracts';
import { getProject, updateContractAddress } from '@store/actions/projects';
import { metamaskError, notify } from '@/internal/notifications';
import { NavbarSettings } from '@views/Home';
import { RootState } from '@/app/store';
import { BadRequestError } from '@/api';
import { ServicePlugin } from '@/plugins/service';
import { switchNetwork } from '@/ethers/network';

import './index.scss';

/** Default sale terms, the same as on the server (nfts.DefaultMintPrice, nfts.DefaultMaxPerMint). */
const DEFAULT_MINT_PRICE: string = '0.001';
const DEFAULT_MAX_PER_MINT: number = 5;
const MIN_PER_MINT: number = 1;
const NO_TOKENS: number = 0;
/** A non-negative amount of ETH the contract can take exactly, in wei. */
const ETHER_AMOUNT: RegExp = /^\d{1,9}(\.\d{1,18})?$/;

const CreateContract: React.FC = () => {
    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();
    const ethersService = ServicePlugin.create();
    const conractSignature = useSelector((state: RootState) => state.contractsReducer.contractSignature);
    const project = useSelector((state: RootState) => state.collectionsReducer.project);
    const navbarSettings: NavbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);
    const [isPublishedPopupVisible, setIsPublishedPopupVisible] = useState<boolean>(false);
    const [name, setName] = useState<string>('');
    const [mintPrice, setMintPrice] = useState<string>(DEFAULT_MINT_PRICE);
    const [maxPerMint, setMaxPerMint] = useState<string>(String(DEFAULT_MAX_PER_MINT));
    // the contract mints exactly as many tokens as the collection has metadata for.
    const nftCount: number = project.id === id ? project.nftCount : NO_TOKENS;

    /** Changes contract name. */
    const changeName = (e: React.ChangeEvent<HTMLInputElement>) => {
        setName(e.target.value);
    };

    /** Changes mint price. */
    const changeMintPrice = (e: React.ChangeEvent<HTMLInputElement>) => {
        setMintPrice(e.target.value.trim());
    };

    /** Changes the limit of tokens per transaction. */
    const changeMaxPerMint = (e: React.ChangeEvent<HTMLInputElement>) => {
        setMaxPerMint(e.target.value);
    };

    /** Returns why the sale terms cannot go into the contract, or an empty string. */
    const termsError = (): string => {
        if (!nftCount) {
            return 'The collection has no tokens';
        }
        if (!ETHER_AMOUNT.test(mintPrice)) {
            return 'Mint price must be a number of ETH, e.g. 0.001';
        }
        const limit = Number(maxPerMint);
        if (!Number.isInteger(limit) || limit < MIN_PER_MINT || limit > nftCount) {
            return `Max per transaction must be from 1 to ${nftCount}`;
        }

        return '';
    };

    /** Requests the deployment data with the sale terms; the effect below sends it from MetaMask. */
    const deployContract = async() => {
        const error = termsError();
        if (error) {
            notify(error);

            return;
        }

        try {
            await switchNetwork();
        } catch (error: any) {
            metamaskError(error);

            return;
        }

        try {
            await dispatch(getContract(id, { name: name.trim(), price: mintPrice, maxPerMint: Number(maxPerMint) }));
        } catch (error: any) {
            notify(error instanceof BadRequestError ? error.message : 'Could not prepare the contract');
        }
    };

    useEffect(() => {
        dispatch(getProject(id));
    }, [id]);

    // the contract name and the limit start from the collection.
    useEffect(() => {
        if (project.id !== id) {
            return;
        }
        setName(name => name || project.name);
        setMaxPerMint(String(Math.min(DEFAULT_MAX_PER_MINT, project.nftCount || DEFAULT_MAX_PER_MINT)));
    }, [project.id]);

    /** Closes publish project`s popup and open preview project page. */
    const closePublishedPopup = () => {
        setIsPublishedPopupVisible(false);
    };

    useEffect(() => {
        conractSignature &&
            (async() => {
                try {
                    const receipt = await ethersService.deployContract(conractSignature);
                    await dispatch(updateContractAddress(id, receipt.contractAddress, receipt.blockNumber));
                    await dispatch(getProject(id));
                    setIsPublishedPopupVisible(true);
                    dispatch(setContract(''));
                } catch (error: any) {
                    metamaskError(error);
                }
            })();
    }, [conractSignature]);

    return <>
        <Navbar navbarSettings={navbarSettings} />
        <ProjectNavigation />
        {
            // TODO: Royalties popup (second contract).
        }
        {
            isPublishedPopupVisible &&
            <Popup
                // TODO: Redirect action to project`s/contract page.
                setIsPopupVisible={setIsPublishedPopupVisible}
            >
                <h1 className="popup__title-middle">
                    Your Contract has been published
                </h1>
                <button
                    className="popup__open"
                    onClick={closePublishedPopup}
                >
                    Confirm
                </button>
            </Popup>
        }
        <section className="create-contract">
            <h1 className="create-contract__title">
                Contract
            </h1>
            <div
                className="create-contract__form"
            >
                <label
                    htmlFor="name"
                    className="create-contract__form__name-label"
                >
                    Name
                </label>
                <input
                    type="text"
                    className="create-contract__form__name"
                    id="name"
                    min={MIN_PER_MINT}
                    value={name}
                    onChange={changeName}
                />
                <StyledSelect />
                <label
                    htmlFor="nft-count"
                    className="create-contract__form__nft-count-label"
                >
                    Number Of NFTs
                </label>
                <input
                    type="number"
                    className="create-contract__form__nft-count"
                    id="nft-count"
                    disabled
                    value={nftCount || ''}
                    readOnly
                />
                <div className="create-contract__form__prices">
                    <div className="create-contract__form__prices__wrapper">
                        <label
                            htmlFor="mint"
                            className="create-contract__form__mint-label"
                        >
                            Mint Price
                        </label>
                        <input
                            type="text"
                            inputMode="decimal"
                            className="create-contract__form__mint"
                            id="mint"
                            value={mintPrice}
                            onChange={changeMintPrice}
                        />
                        <span className="create-contract__form__mint-placeholder">
                            ETH
                        </span>
                    </div>
                    <div className="create-contract__form__prices__wrapper">
                        <label
                            htmlFor="max-per-mint"
                            className="create-contract__form__whitelist-label"
                        >
                            Max Per Transaction
                        </label>
                        <input
                            type="number"
                            className="create-contract__form__whitelist"
                            id="max-per-mint"
                            min={MIN_PER_MINT}
                            max={nftCount || undefined}
                            value={maxPerMint}
                            onChange={changeMaxPerMint}
                        />
                    </div>
                </div>
                <NavLink
                    className="create-contract__form__whitelist-btn"
                    to={`/whitelist/${id}`}
                >
                    <span className="create-contract__form__whitelist-btn__text">Whitelist</span>
                </NavLink>
                <button
                    className={`create-contract__form__submit${name ? '' : '-disabled'}`}
                    disabled={!name}
                    onClick={deployContract}
                >Publish</button>
            </div>
        </section>
    </>;
};

export default CreateContract;
