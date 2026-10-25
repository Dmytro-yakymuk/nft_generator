import { DashboardClient } from '@/api/dashboard';
import { CollectionEventsPage, CollectionStatistics, SalesDashboard } from '@/dashboard';

/**
 * Exposes all dashboard related logic.
 */
export class DashboardService {
    private readonly dashboard: DashboardClient;

    /** DashboardService contains http implementation of dashboard API. */
    public constructor(dashboard: DashboardClient) {
        this.dashboard = dashboard;
    };

    /** Gets a page of the collection events. */
    public async listCollectionEvents(projectId: string, page: string): Promise<CollectionEventsPage> {
        return await this.dashboard.listCollectionEvents(projectId, page);
    }

    /** Gets dashboard collection statistics. */
    public async getCollectionStatistics(projectId: string): Promise<CollectionStatistics> {
        return await this.dashboard.getCollectionStatistics(projectId);
    }

    /** Gets paid events of the period and their prices. */
    public async getSalesDashboard(projectId: string, period: number, from: number, to: number): Promise<SalesDashboard> {
        return await this.dashboard.getSalesDashboard(projectId, period, from, to);
    }
}
