import { useSelector } from 'react-redux';

import { toFixedNoRounding } from '@/app/internal/balance';
import { RootState } from '@/app/store';

import './index.scss';

/** Sales price value entity. */
class SalesPrice {
    public constructor(
        public title: string,
        public price: number | string,
    ) {}
}

export const ProjectValues: React.FC= () => {
    const NUMBER_OF_CHARACTERS_AFTER_DOT: number = 8;
    const salesPrices = useSelector((state: RootState) => state.dashboardReducer.salesDashboard.prices);

    /** List of sales prices of the chart period for the dashboard. */
    const dashboardPrices: SalesPrice[] = [
        new SalesPrice('Floor Price', salesPrices.floorPrice),
        new SalesPrice('Average Price', salesPrices.averagePrice),
        new SalesPrice('Highest Sale', salesPrices.highestPrice),
        new SalesPrice('Volume', salesPrices.volume),
    ];

    return (
        <div className="dashboard__values">
            {dashboardPrices.map((item: SalesPrice, index: number) =>
                <div className="dashboard__values__item" key={index}>
                    <span className="price">{toFixedNoRounding(+item.price, NUMBER_OF_CHARACTERS_AFTER_DOT)} ETH</span>
                    <span className="title">{item.title}</span>
                </div>
            )}
        </div>
    );
};
