import { SET_COLLECTION_STATISTICS, SET_EVENTS, SET_SALES_DASHBOARD } from '@store/actions/dashboard';
import { CollectionEventsPage, CollectionStatistics, SalesDashboard, SalesPrices } from '@/dashboard';

/** Exposes dashboard entity store state. Uses as default reducer state. */
class DashboardState {
    constructor(
        public eventsPage: CollectionEventsPage = new CollectionEventsPage([], '', ''),
        public collectionStatistics: CollectionStatistics = new CollectionStatistics(),
        public salesDashboard: SalesDashboard = new SalesDashboard(
            [],
            new SalesPrices()
        ),
    ) { };
};

/** Exposes action for dashboard reducer. */
class DashboardReducerAction {
    /** At least, contains action type and payload data. */
    constructor(
        public type: string = '',
        public payload: any = '',
    ) { };
};

export const dashboardReducer = (
    dashboardState: DashboardState = new DashboardState(),
    action: DashboardReducerAction = new DashboardReducerAction(),
) => {
    switch (action.type) {
    case SET_EVENTS:
        dashboardState.eventsPage = action.payload;
        break;
    case SET_COLLECTION_STATISTICS:
        dashboardState.collectionStatistics = action.payload;
        break;
    case SET_SALES_DASHBOARD:
        dashboardState.salesDashboard = action.payload;
        break;
    default:
        break;
    };

    return { ...dashboardState };
};
