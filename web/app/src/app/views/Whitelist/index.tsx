import { useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import { Navbar } from '@components/common/Navbar';
import { Paginator } from '@components/common/Paginator';
import { Popup } from '@components/common/Popup';
import { ProjectNavigation } from '@components/common/ProjectNavigation';
import { ContractMerkleRoot } from '@components/whitelist/ContractMerkleRoot';
import { NavbarSettings } from '@views/Home';

import { BadRequestError, InternalError } from '@/api';
import { useDebounce } from '@/app/hooks/useDebounce';
import { RootState } from '@/app/store';
import { addWalletToWhitelist, deleteWalletFromWhitelist, getWhitelist, uploadWhitelistFile } from '@/app/store/actions/whitelist';
import { DEFAULT_CURRENT_PAGE } from '@/app/store/reducers/whitelist';
import { convertToBase64 } from '@/internal/convertImage';
import { notify } from '@/internal/notifications';
import { WhitelistedWallet } from '@/whitelist';

import { Close } from '@static/images/dashboard/index';
import pinkLogo from '@static/images/logo.png';
import importIcon from '@static/images/import.svg';
import search from '@static/images/search.svg';
import file from '@static/images/whitelist/file.svg';
import previewFile from '@static/images/whitelist/preview-file.svg';
import addFile from '@static/images/whitelist/add-file.svg';
import cross from '@static/images/whitelist/cross.svg';

import './index.scss';

const Whitelist: React.FC = () => {
    const FIRST_LIST_INDEX = 0;
    const SECOND_LIST_INDEX = 1;
    /** Debounce delay in miliseconds. */
    const DEBONCE_DELAY: number = 500;

    const dispatch = useDispatch();
    const xlsxRef = useRef<HTMLInputElement | null>(null);
    const labelRef = useRef<HTMLLabelElement | null>(null);
    const { id }: { id: string } = useParams();
    const navbarSettings: NavbarSettings = new NavbarSettings('border-pink', 'text-black', 'background-white', pinkLogo);
    const { wallets, page } = useSelector((state: RootState) => state.whitelistReducer.whitelistPage);
    const [walletAddress, setWalletAddress] = useState<string>('');
    const [walletName, setWalletName] = useState<string>('');
    const [searchedWalletName, setSearchedWalletName] = useState<string>('');
    const [isAddToWhitelistVisible, setIsAddToWhitelistVisible] = useState<boolean>(false);
    const [isUploadFilePopupVisible, setIsUploadFilePopupVisible] = useState<boolean>(false);
    const [isFileSelected, setIsFileSelected] = useState<boolean>(false);
    const [selectedBase64String, setSelectedBase64String] = useState<string>('');
    const selectedFileClassName = isFileSelected ? '' : '-not-selected';

    const [currentPage, setCurrentPage] = useState<number>(page.currentPage);
    const debouncedValue: string = useDebounce(searchedWalletName, DEBONCE_DELAY);

    /** Handle open whitelist popup. */
    const handleOpenWhiteListPopup = () => {
        setIsAddToWhitelistVisible(true);
    };

    /** Handle close Whitelist popup. */
    const handleClosePopup = () => {
        setIsAddToWhitelistVisible(false);
    };

    /** Handle add wallet address to whitelist. */
    const handleAddWalletToWhitelist = async() => {
        const walletPattern = /^0x[0-9a-f]{40}$/i;

        if (!walletName) {
            notify('Enter wallet name!');

            return;
        }

        if (!walletPattern.test(walletAddress)) {
            notify('Invalid wallet address!');

            return;
        }

        try {
            await dispatch(addWalletToWhitelist(id, walletName, walletAddress));
            setWalletAddress('');
            setWalletName('');
            setIsAddToWhitelistVisible(false);
        } catch (error) {
            if (error instanceof InternalError) {
                notify('This wallet address already exists!');

                return;
            }
            notify('Something went wrong');
        }
    };

    /** Handle uploads xlxs whitelist file, parse to base64 string and sends to the server. */
    const handleUploadExcelFile = () => {
        const inputFile = xlsxRef.current;
        const labelFile = labelRef.current;

        if (!inputFile || !labelFile) {
            return;
        }

        inputFile && labelFile &&
            inputFile.addEventListener('change', async() => {
                inputFile.files &&
                    (labelFile.innerText = inputFile.files[FIRST_LIST_INDEX].name);
                setIsFileSelected(true);
                const convertedXlxsToBase64 = inputFile.files && await convertToBase64(inputFile.files[FIRST_LIST_INDEX]);
                const base64String = convertedXlxsToBase64.split('sheet;base64,')[SECOND_LIST_INDEX];
                setSelectedBase64String(base64String);
            });
    };

    /** Sends selected xlxs file to the server. */
    const sendUploadFileToServer = async() => {
        try {
            await dispatch(uploadWhitelistFile(id, selectedBase64String));
            setIsFileSelected(false);
            setIsUploadFilePopupVisible(false);
            notify('File uploaded successfully', 'success');
        } catch (error) {
            if (error instanceof BadRequestError) {
                notify('Addresses already exists.');

                return;
            }
            notify('Something wrong. Please try later.');
        }
    };

    /** Deletes selected file from input. */
    const deleteSelectedFile = () => {
        const inputFile = xlsxRef.current;

        if (!inputFile) {
            return;
        }

        inputFile.value = '';
        setIsFileSelected(false);
    };

    /** Deletes wallet from whitelist. */
    const deleteWhitelistWallet = (address: string) => {
        dispatch(deleteWalletFromWhitelist(id, address, currentPage));
    };

    /** Changes wallet name. */
    const changeWalletName = (e: React.ChangeEvent<HTMLInputElement>) => {
        setWalletName(e.target.value);
    };

    /** Changes searched wallet name. */
    const changeSearchedWalletName = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearchedWalletName(e.target.value);
        setCurrentPage(DEFAULT_CURRENT_PAGE);
    };

    /** Changes searched wallet address. */
    const changeWalletAddress = (e: React.ChangeEvent<HTMLInputElement>) => {
        setWalletAddress(e.target.value);
    };

    useEffect(() => {
        dispatch(getWhitelist(id, currentPage, debouncedValue));
    }, [currentPage, debouncedValue]);

    return (
        <>
            <Navbar navbarSettings={navbarSettings} />
            <ProjectNavigation />
            <div className="dashboard__white-list">
                {
                    isAddToWhitelistVisible &&
                    <Popup
                        setIsPopupVisible={setIsAddToWhitelistVisible}
                    >
                        <h1 className="popup__title-left-side">
                            Add to Whitelist
                        </h1>
                        <label htmlFor="wallet-name">Name</label>
                        <input
                            id="wallet-name"
                            type="text"
                            placeholder="Wallet name"
                            onChange={changeWalletName}
                            value={walletName}
                        />
                        <label htmlFor="wallet-address">Wallet Address</label>
                        <input
                            id="wallet-address"
                            type="text"
                            onChange={changeWalletAddress}
                            value={walletAddress}
                        />
                        <div className="popup__buttons">
                            <button
                                className="popup__cancel"
                                onClick={handleClosePopup}
                            >
                                Cancel
                            </button>
                            <button
                                className="popup__confirm"
                                onClick={handleAddWalletToWhitelist}
                            >
                                Create
                            </button>
                        </div>
                    </Popup>
                }
                {
                    isUploadFilePopupVisible &&
                    <Popup
                        setIsPopupVisible={setIsUploadFilePopupVisible}
                    >
                        <h1 className="upload-popup__title">Select an XLXS File to upload</h1>
                        <div className="upload-popup__file">
                            <img className="upload-popup__file__icon" src={file} alt="File" />
                            <img className="upload-popup__file__add-icon" src={addFile} alt="Add file" />
                            <label
                                htmlFor="whitelist"
                                className="upload-popup__file__label"
                                onClick={handleUploadExcelFile}
                            ></label>
                            <span className="upload-popup__file__text">Drag your file here or click to browse.</span>
                            <input
                                type="file"
                                accept=".xlsx"
                                id="whitelist"
                                ref={xlsxRef}
                                hidden
                                name="whitelist"
                                onChange={() => {}}
                            />
                        </div>
                        <div className={`upload-popup__preview${selectedFileClassName}`}>
                            <img className="upload-popup__preview__icon" src={previewFile} alt="Preview" />
                            <span className="upload-popup__preview__name" ref={labelRef}></span>
                            <div className="upload-popup__preview__delete" onClick={deleteSelectedFile}>
                                <img className="upload-popup__preview__delete__icon" src={cross} alt="Delete" />
                            </div>
                        </div>
                        <button
                            className={`upload-popup__btn${selectedFileClassName}`}
                            disabled={!isFileSelected}
                            onClick={sendUploadFileToServer}
                        >Upload</button>
                    </Popup>
                }
                <div className="dashboard__white-list__wrapper">
                    <div className="dashboard__white-list__header">
                        <span className="dashboard__white-list__title">Whitelist</span>
                        <div className="dashboard__white-list__btn import" onClick={() => setIsUploadFilePopupVisible(true)}>
                            <label className="dashboard__white-list__btn-text">Import</label>
                            <img src={importIcon} alt="Import" />
                        </div>
                        <button className="dashboard__white-list__btn add" onClick={handleOpenWhiteListPopup}>
                            <span className="dashboard__white-list__btn-text">Add to Whitelist</span>
                        </button>
                    </div>
                    <ContractMerkleRoot projectId={id} wallets={wallets} />
                    <div className="dashboard__white-list__search">
                        <input
                            type="text"
                            placeholder="Search..."
                            value={searchedWalletName}
                            onChange={changeSearchedWalletName}
                        />
                        <img src={search} alt="Search" />
                    </div>
                    <div className="dashboard__white-list__divider"></div>
                    <div className="dashboard__white-list__websites">
                        {
                            wallets &&
                            wallets.map(
                                (wallet: WhitelistedWallet, index: number) =>
                                    <div key={index} className="dashboard__white-list__website">
                                        <span className="dashboard__white-list__website__name">{wallet.name}</span>
                                        <span className="dashboard__white-list__website__address">{wallet.walletAddress}</span>
                                        <div
                                            onClick={() => deleteWhitelistWallet(wallet.walletAddress)}
                                            className="dashboard__white-list__website__close"
                                        >
                                            <Close />
                                        </div>
                                    </div>
                            )
                        }
                    </div>
                </div>
            </div>
            <Paginator
                currentPage={currentPage}
                setCurrentPage={setCurrentPage}
                itemsCount={page.totalCount}
            />
        </>
    );
};

export default Whitelist;
