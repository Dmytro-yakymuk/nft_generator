import { APIClient } from '@/api';
import { CollectionEvent, CollectionEventsPage, CollectionStatistics, SalesDashboard, SalesPrices } from '@/dashboard';

const ZERO_DATE_PERIOD = 0;

/** Creates a collection event from the response item. */
const toCollectionEvent = (event: CollectionEvent) =>
    new CollectionEvent(
        event.from,
        event.to,
        event.tokenId,
        event.timestamp,
        event.type,
        event.price,
        event.currency,
        event.image,
        event.txHash,
    );

/**
 * DashboardClient is a http implementation of dashboard API.
 * The server builds the analytics from the events of the collection contract in the blockchain.
 */
export class DashboardClient extends APIClient {
    private readonly ROOT_PATH: string = '/api/v0/dashboard';

    /** Gets a page of the collection events, newest first. */
    public async listCollectionEvents(projectId: string, page: string): Promise<CollectionEventsPage> {
        const path = `${this.ROOT_PATH}/${projectId}/events${page && `?page=${page}`}`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const eventsPage = await response.json();

        return new CollectionEventsPage(
            eventsPage.events.map(toCollectionEvent),
            eventsPage.nextPage,
            eventsPage.previousPage
        );
    };

    /** Gets dashboard collection statistics. */
    public async getCollectionStatistics(projectId: string): Promise<CollectionStatistics> {
        const path = `${this.ROOT_PATH}/${projectId}/stats`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const statistics = await response.json();

        return new CollectionStatistics(statistics.items, statistics.owners, statistics.sales);
    };

    /** Gets paid events of the period and their prices. */
    public async getSalesDashboard(projectId: string, period: number, from: number, to: number): Promise<SalesDashboard> {
        const daysPeriod = period === ZERO_DATE_PERIOD ? '' : period;
        const unixDateFrom = from === ZERO_DATE_PERIOD ? '' : from;
        const unixDateTo = to === ZERO_DATE_PERIOD ? '' : to;
        const path = `${this.ROOT_PATH}/${projectId}/sales?period=${daysPeriod}&from=${unixDateFrom}&to=${unixDateTo}`;
        const response = await this.http.get(path);

        if (!response.ok) {
            await this.handleError(response);
        }

        const dashboard = await response.json();

        return new SalesDashboard(
            dashboard.events.map(toCollectionEvent),
            new SalesPrices(
                dashboard.floorPrice,
                dashboard.highestPrice,
                dashboard.averagePrice,
                dashboard.volume,
            )
        );
    };
};
