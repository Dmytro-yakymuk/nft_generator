import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Pagination } from '@components/common/Pagination';

import { RootState } from '@/app/store';
import { listCollectionEvents } from '@/app/store/actions/dashboard';
import { notify } from '@/internal/notifications';
import { CollectionEvent } from '@/dashboard';
import { toFixedNoRounding } from '@/app/internal/balance';
import { convertDateTime } from '@/app/internal/convertDate';
import { getNetwork } from '@/ethers/network';

import './index.scss';

const FIRST_ADDRESS_INDEX = 0;
const SECOND_ADDRESS_INDEX = 6;
const THIRD_ADDRESS_INDEX = 36;
const FOURTH_ADDRESS_INDEX = 42;
const MESSAGE_TIME = 2000;
const NUMBER_OF_CHARACTERS_AFTER_DOT: number = 5;

export const Activity: React.FC<{ projectId: string }> = ({ projectId }) => {
    const [isFullAddressVisible, setIsFullAddressVisible] = useState<boolean>(false);
    const [isLoaded, setIsLoaded] = useState<boolean>(false);
    const dispatch = useDispatch();
    const eventsPage = useSelector((state: RootState) => state.dashboardReducer.eventsPage);
    const { blockExplorerUrl } = getNetwork();

    useEffect(() => {
        (async() => {
            try {
                await dispatch(listCollectionEvents(projectId, ''));
                setIsLoaded(true);
            } catch (error) {
                notify('Could not read the collection contract from the network. Try later.');
            }
        })();
    }, [projectId]);

    /** Short contract address view. */
    const changeAddressVisability = (address: string): string => {
        const firstPart = address.slice(FIRST_ADDRESS_INDEX, SECOND_ADDRESS_INDEX);
        const lastPart = address.slice(THIRD_ADDRESS_INDEX, FOURTH_ADDRESS_INDEX);

        if (isFullAddressVisible) {
            return address;
        }

        return `${firstPart}...${lastPart}`;
    };

    /** Copies address and shows copy message. */
    const alertCopiedMessage = (address: string) => {
        const copiedBlock = document.querySelector<HTMLDivElement>('.copied-block');

        if (!copiedBlock) { return; }
        copiedBlock.setAttribute('style', 'display: flex;');

        const copiedText: HTMLTextAreaElement = document.createElement('textarea');
        document.body.appendChild(copiedText);

        copiedText.value = address;

        copiedText.select();
        document.execCommand('copy');
        document.body.removeChild(copiedText);

        setTimeout(() => {
            if (!copiedBlock) { return; }
            copiedBlock.setAttribute('style', 'display: none;');
        }, MESSAGE_TIME);
    };

    return (
        <>
            <div className="copied-block">Address copied</div>
            <table className="dashboard__activity">
                <thead className="dashboard__activity__head">
                    <tr className="dashboard__activity__titles">
                        <th className="dashboard__activity__caption">Activity</th>
                        <th className="dashboard__activity__caption"></th>
                        <th className="dashboard__activity__title">From</th>
                        <th className="dashboard__activity__title">To</th>
                        <th className="dashboard__activity__title">Date</th>
                        <th className="dashboard__activity__title">Price</th>
                    </tr>
                </thead>
                <tbody className="dashboard__activity__body">
                    {isLoaded && !eventsPage.events.length &&
                        <tr className="dashboard__activity__values">
                            <td className="dashboard__activity__value" colSpan={6}>No mints or transfers yet</td>
                        </tr>
                    }
                    {eventsPage.events.map((event: CollectionEvent, index: number) =>
                        <tr className="dashboard__activity__values" key={index}>
                            <td className="dashboard__activity__value">
                                {event.image && <img src={event.image} alt={`Token #${event.tokenId}`} />}
                            </td>
                            <td className="dashboard__activity__value">{event.type} #{event.tokenId}</td>
                            <td
                                className="dashboard__activity__value"
                                onClick={() => alertCopiedMessage(event.from)}
                            >{changeAddressVisability(event.from)}</td>
                            <td
                                className="dashboard__activity__value"
                                onClick={() => alertCopiedMessage(event.to)}
                            >{changeAddressVisability(event.to)}</td>
                            <td className="dashboard__activity__value">
                                {blockExplorerUrl
                                    ? <a href={`${blockExplorerUrl}/tx/${event.txHash}`} target="_blank" rel="noreferrer">
                                        {convertDateTime(event.timestamp)}
                                    </a>
                                    : convertDateTime(event.timestamp)}
                            </td>
                            <td className="dashboard__activity__price">
                                {+event.price ? `${toFixedNoRounding(+event.price, NUMBER_OF_CHARACTERS_AFTER_DOT)} ${event.currency}` : '—'}
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
            <Pagination projectId={projectId} />
        </>
    );
};
