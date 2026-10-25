import { useDispatch, useSelector } from 'react-redux';

import { NextIcon, PreviousIcon } from '@static/images/dashboard';

import { RootState } from '@/app/store';
import { listCollectionEvents } from '@/app/store/actions/dashboard';
import { notify } from '@/internal/notifications';

import './index.scss';

export const Pagination: React.FC<{ projectId: string }> = ({ projectId }) => {
    const dispatch = useDispatch();
    const eventsPage = useSelector((state: RootState) => state.dashboardReducer.eventsPage);

    /** Goes to the next events page. */
    const goToNextPage = async() => {
        if (eventsPage.nextPage) {
            try {
                await dispatch(listCollectionEvents(projectId, eventsPage.nextPage));
            } catch (error) {
                notify('Something wrong. Try later.');
            }
        }
    };

    /** Goes to the previous events page. */
    const goToPrevPage = async() => {
        if (eventsPage.previousPage) {
            try {
                await dispatch(listCollectionEvents(projectId, eventsPage.previousPage));
            } catch (error) {
                notify('Something wrong. Try later.');
            }
        }
    };

    return (
        <div className="activity-pagination">
            <div
                className={`activity-pagination__prev${eventsPage.previousPage ? '' : '-disabled'}`}
                onClick={goToPrevPage}
            >
                <PreviousIcon />
                <span className="activity-pagination__btn">Prev</span>
            </div>
            <div className={`activity-pagination__next${eventsPage.nextPage ? '' : '-disabled'}`}
                onClick={goToNextPage}
            >
                <span className="activity-pagination__btn">Next</span>
                <NextIcon />
            </div>
        </div>
    );
};
