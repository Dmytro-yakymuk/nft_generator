import { WhitelistPage } from '@/whitelist';
import { DELETE_FROM_WHITELIST, SET_ADDRESS_TO_WHITELIST, SET_WHITELIST } from '@store/actions/whitelist';
import { ActionType } from '@/app/store';

/** Default pagination parameters.
 * Uses as mock page values for initialize redux.
 */
const DEFAULT_OFFSET_VALUE: number = 0;
const DEFAULT_LIMIT_VALUE: number = 10;
export const DEFAULT_CURRENT_PAGE: number = 1;
const PAGES_COUNT: number = 1;
const ITEMS_TOTAL_COUNT: number = 1;

const page = {
    offset: DEFAULT_OFFSET_VALUE,
    limit: DEFAULT_LIMIT_VALUE,
    currentPage: DEFAULT_CURRENT_PAGE,
    pageCount: PAGES_COUNT,
    totalCount: ITEMS_TOTAL_COUNT,
};

const defaultWhitelistPage = new WhitelistPage([], page);

/** Exposes whitelist entity store state. Uses as default reducer state. */
class WhitelistState {
    constructor(
        public whitelistPage: WhitelistPage = defaultWhitelistPage,
    ) { };
};

export const whitelistReducer = (
    whitelistState: WhitelistState = new WhitelistState(),
    action: ActionType<any> = new ActionType('', ''),
) => {
    switch (action.type) {
    case SET_WHITELIST:
        whitelistState.whitelistPage = action.payload;
        break;
    case SET_ADDRESS_TO_WHITELIST:
        whitelistState.whitelistPage.wallets = [...whitelistState.whitelistPage.wallets, action.payload];
        break;
    case DELETE_FROM_WHITELIST:
        whitelistState.whitelistPage.wallets = action.payload;
        break;
    default:
        break;
    };

    return { ...whitelistState };
};
