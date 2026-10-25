import { useState } from 'react';

import ethereum from '@static/images/ethereum.png';
import flow from '@static/images/flow.png';
import polygon from '@static/images/polygon.png';
import solana from '@static/images/solana.png';

import './index.scss';

/** Exposes available network entity with options. */
class NetworkOptions {
    /** NetworkOptions at least, contains image, name and status of network. */
    constructor(
        public image: string,
        public name: string,
        public status: string,
    ) { };
};

/** StyledSelect is custom select component. */
const StyledSelect: React.FC = () => {
    /** Index of the first sort item. */
    const FIRST_SORT_ITEMS_INDEX: number = 0;

    /** Describes available networks. */
    const networkOptionsConfig: NetworkOptions[] = [
        new NetworkOptions(ethereum, 'Ethereum', ''),
        new NetworkOptions(flow, 'Flow', 'Coming soon'),
        new NetworkOptions(polygon, 'Polygon', 'Coming soon'),
        new NetworkOptions(solana, 'Solana', 'Coming soon'),
    ];

    /** Indicates active network options. */
    const [activeNetworkOptions, setActiveNetworkOptions] = useState<NetworkOptions>(
        networkOptionsConfig[FIRST_SORT_ITEMS_INDEX]
    );

    /** Indicates if networks shown. */
    const [isNetworksShown, setIsNetworksShown] = useState<boolean>(false);

    /** Changes netwoks visavility. */
    const changeNetworksVisability = () => {
        setIsNetworksShown(!isNetworksShown);
    };

    return (
        <div
            className="styled-select"
            onClick={changeNetworksVisability}
        >
            <div className="styled-select__selected-item">
                <img
                    src={activeNetworkOptions.image}
                />
                <span className="styled-select__selected-item__name">
                    {activeNetworkOptions.name}
                </span>
                <span className="styled-select__selected-item__status">
                    {activeNetworkOptions.status}
                </span>
            </div>
            <div className={`styled-select__items${isNetworksShown ? '-active' : ''}`} >
                <ul>
                    {networkOptionsConfig.map((networOptions: NetworkOptions, index: number) =>
                        <li
                            className="styled-select__items-active__item"
                            key={index}
                        >
                            <img src={networOptions.image} />
                            <span className="styled-select__items-active__item__name">
                                {networOptions.name}
                            </span>
                            <span className="styled-select__items-active__item__status">
                                {networOptions.status}
                            </span>
                        </li>
                    )}
                </ul>
            </div>
        </div>
    );
};

export default StyledSelect;
