import { Dispatch } from 'redux';

import { DashboardClient } from '@/api/dashboard';
import { DashboardService } from '@/dashboard/service';
import { CollectionEventsPage, CollectionStatistics, SalesDashboard } from '@/dashboard';

export const SET_EVENTS = 'SET_EVENTS';
export const SET_COLLECTION_STATISTICS = 'SET_COLLECTION_STATISTICS';
export const SET_SALES_DASHBOARD = 'SET_SALES_DASHBOARD';

/** An action setCollectionEvents contains type and payload data for dashboard events. */
export const setCollectionEvents = (eventsPage: CollectionEventsPage) => ({
    type: SET_EVENTS,
    payload: eventsPage,
});

/** An action setCollectionStatistics contains type and payload data for dashboard collection statistics. */
export const setCollectionStatistics = (statistics: CollectionStatistics) => ({
    type: SET_COLLECTION_STATISTICS,
    payload: statistics,
});

/** An action setSalesDashboard contains type and payload data for dashboard sales chart and prices. */
export const setSalesDashboard = (salesDashboard: SalesDashboard) => ({
    type: SET_SALES_DASHBOARD,
    payload: salesDashboard,
});

const dashboardClient = new DashboardClient();
const dashboardService = new DashboardService(dashboardClient);

/** Thunk middleware that gets a page of the collection events from api and sets it into state */
export const listCollectionEvents = (projectId: string, page: string) => async function(dispatch: Dispatch) {
    const eventsPage = await dashboardService.listCollectionEvents(projectId, page);
    dispatch(setCollectionEvents(eventsPage));
};

/** Thunk middleware that gets collection statistics from api and sets it into state */
export const getCollectionStatistics = (projectId: string) => async function(dispatch: Dispatch) {
    const collectionStatistics = await dashboardService.getCollectionStatistics(projectId);
    dispatch(setCollectionStatistics(collectionStatistics));
};

/** Thunk middleware that gets paid events of the period and their prices from api and sets it into state */
export const getSalesDashboard = (projectId: string, period: number, from: number, to: number) => async function(dispatch: Dispatch) {
    const salesDashboard = await dashboardService.getSalesDashboard(projectId, period, from, to);
    dispatch(setSalesDashboard(salesDashboard));
};
