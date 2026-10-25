import { Dispatch, SetStateAction } from 'react';

import './index.scss';

/** Exposes paginator props.
 * Uses as props types for common paginator component.
 */
export class PaginatorProps {
    constructor(
        public currentPage: number,
        public itemsCount: number,
        public setCurrentPage: Dispatch<SetStateAction<number>>,
    ) {}
};

/** Commont paginator component. */
export const Paginator: React.FC<PaginatorProps> = ({ currentPage, itemsCount, setCurrentPage }) => {
    /** Default selected page. */
    const DEFAULT_SELECTED_PAGE: number = 1;
    /** Default item number. */
    const DEFAULT_ITEM_NUMBER: number = 1;
    /** Describes default items on page count. */
    const ITEMS_ON_PAGE: number = 10;

    const pagesCount = Math.ceil(itemsCount / ITEMS_ON_PAGE);

    /** Indicates if first page is selected. */
    const isFirstPageSelected: boolean = currentPage === DEFAULT_SELECTED_PAGE;
    /** Indicates if last page is selected. */
    const isLastPageSelected: boolean = currentPage === pagesCount || pagesCount === itemsCount;

    /** Lower shown items breakpoint. */
    const LOWER_SHOWN_ITEMS_VALUE_BREAKPOINT: number = isFirstPageSelected ?
        DEFAULT_ITEM_NUMBER : (currentPage - DEFAULT_SELECTED_PAGE) * ITEMS_ON_PAGE + DEFAULT_ITEM_NUMBER;
    /** Upper shown items breakpoint. */
    const UPPER_SHOWN_ITEMS_VALUE_BREAKPOINT: number = isLastPageSelected ? itemsCount : currentPage * ITEMS_ON_PAGE;

    /** Changes previous page. */
    const changePrevPage = () => {
        if (currentPage <= DEFAULT_SELECTED_PAGE) {
            return;
        }
        setCurrentPage(currentPage => currentPage - DEFAULT_SELECTED_PAGE);
    };

    /** Changes next page. */
    const changeNextPage = () => {
        if (currentPage >= pagesCount) {
            return;
        }
        setCurrentPage(currentPage => currentPage + DEFAULT_SELECTED_PAGE);
    };

    /** Logic for showing paging result. */
    const paginatorVisibilityLogic = !itemsCount
        ? `0-0 of ${itemsCount}`
        :`${LOWER_SHOWN_ITEMS_VALUE_BREAKPOINT}-${UPPER_SHOWN_ITEMS_VALUE_BREAKPOINT} of ${itemsCount}`;

    return <section className="paginator">
        <div className="paginator__wrapper">
            <div className="paginator__all-count">
                {paginatorVisibilityLogic}
            </div>
            <div className="paginator__sellect-page">
                <button
                    className={`paginator__sellect-page__prev${isFirstPageSelected ? '-not-active' : ''}`}
                    onClick={changePrevPage}
                >
                    &#8249; Prev
                </button>
                <button
                    className={`paginator__sellect-page__next${isLastPageSelected ? '-not-active' : ''}`}
                    onClick={changeNextPage}
                >
                    Next &#8250;
                </button>
            </div>
        </div>
    </section>;
};
