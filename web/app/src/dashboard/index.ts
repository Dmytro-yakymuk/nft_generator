const DEFAULT_VALUE = 0;

/** Exposes CollectionStatistics domain entity. */
export class CollectionStatistics {
    public constructor(
        public items: number = DEFAULT_VALUE,
        public owners: number = DEFAULT_VALUE,
        public sales: number = DEFAULT_VALUE,
    ) { };
};

/** Exposes collection event entity: a transfer of a collection token read from the blockchain. */
export class CollectionEvent {
    public constructor(
        public from: string,
        public to: string,
        public tokenId: string,
        public timestamp: string,
        /** One of mint, sale, transfer, burn. */
        public type: string,
        /** Ether paid for the token, '0' if the transfer was not paid. */
        public price: string,
        public currency: string,
        public image: string,
        public txHash: string,
    ) {}
}

/** Exposes collection events page entity. */
export class CollectionEventsPage {
    public constructor(
        public events: CollectionEvent[],
        public nextPage: string,
        public previousPage: string,
    ) {}
}

/** Exposes sales prices and volume. */
export class SalesPrices {
    public constructor(
        public floorPrice: string = '',
        public highestPrice: string = '',
        public averagePrice: string = '',
        public volume: string = '',
    ) {}
}

/** Exposes sales dashboard entity, contains paid events of the period and their prices. */
export class SalesDashboard {
    public constructor(
        public events: CollectionEvent[],
        public prices: SalesPrices,
    ) {}
}

/** Chart period entity. */
export class ChartPeriod {
    public constructor(
        public label: string,
        public period: number,
    ) {}
}
