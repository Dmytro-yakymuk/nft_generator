import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
    Area,
    AreaChart,
    LineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from 'recharts';
import { DateRange } from 'react-date-range';

import { CustomTooltip } from './CustomTooltip';
import { Filters } from './Filters';

import { notify } from '@/internal/notifications';
import { RootState } from '@/app/store';
import { convertDateTime } from '@/app/internal/convertDate';
import { calcDateDifference } from '@/app/internal/dateDifference';
import { convertDateToUnix } from '@/app/internal/unixDateFormat';
import { ChartPeriod, CollectionEvent } from '@/dashboard';
import { getSalesDashboard } from '@/app/store/actions/dashboard';

import calendar from '@static/images/dashboard/calendar.svg';
import { Loading } from '@static/images/dashboard/index';

import './index.scss';
import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';

const ONE_DAY_IN_MILLISECONDS = 86400;
/** Prices start from zero; the highest one is drawn below the top of the chart. */
const Y_AXIS_HEADROOM = 1.25;
const CHART_STYLES_NUMBER = 0;
const FIRST_ARRAY_INDEX = 0;
const AXIOS_PADDING = 10;
const YEAR_TO_DATE = 365;
const THREE_MONTHS = 90;
const DOTTER_RADIUS = 7;
const LINE_STROKE = 3;
const MONTH = 30;
const WEEK = 7;
const DAY = 1;
/** Chart period filters. */
const chartPeriods: ChartPeriod[] = [
    new ChartPeriod('All Time', FIRST_ARRAY_INDEX),
    new ChartPeriod('Year To Date', YEAR_TO_DATE),
    new ChartPeriod('3 months', THREE_MONTHS),
    new ChartPeriod('1 months', MONTH),
    new ChartPeriod('7 days', WEEK),
    new ChartPeriod('1 day', DAY),
];

/** Calendar date range pick entity. */
export class DateRangePick {
    public constructor(
        public startDate: Date | null = null,
        public endDate: Date | null = null,
        public key: string = 'selection',
    ) {}
}

export const ProjectChart: React.FC<{ projectId: string }> = ({ projectId }) => {
    const dispatch = useDispatch();
    const salesEvents = useSelector((state: RootState) => state.dashboardReducer.salesDashboard.events);
    const [activeDropdownPeriod, setActiveDropdownPeriod] = useState<ChartPeriod>(chartPeriods[FIRST_ARRAY_INDEX]);
    const [isCalendarVisible, setIsCalendarVisible] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [calendarDateRange, setCalendarDateRange] = useState<any>([new DateRangePick()]);
    const elementClassName: string = !isCalendarVisible ? 'hidden' : '';
    const firstCalendarDateRangeElement = calendarDateRange[FIRST_ARRAY_INDEX];
    const calendarBlockClassName: string = firstCalendarDateRangeElement.startDate ? 'active' : 'inactive';
    /** Selects the date and price from the list of events for the chart. */
    const dashboardData = salesEvents.map((event: CollectionEvent) => ({ 'date': convertDateTime(event.timestamp), 'value': +event.price }));
    const calendarPeriod: number = calcDateDifference(firstCalendarDateRangeElement.startDate, firstCalendarDateRangeElement.endDate);
    /** Calendar returns the number of selected days one less (so here we add one day). */
    const daysDeclension: string = calendarPeriod + DAY === DAY ? 'day' : 'days';
    const unixDateFrom = firstCalendarDateRangeElement.startDate
        ? convertDateToUnix(firstCalendarDateRangeElement.startDate) + ONE_DAY_IN_MILLISECONDS
        : FIRST_ARRAY_INDEX;
    const unixDateTo = firstCalendarDateRangeElement.endDate
        ? convertDateToUnix(firstCalendarDateRangeElement.endDate) + ONE_DAY_IN_MILLISECONDS
        : FIRST_ARRAY_INDEX;

    useEffect(() => {
        (async() => {
            try {
                setIsLoading(true);
                await dispatch(getSalesDashboard(projectId, activeDropdownPeriod.period, unixDateFrom, unixDateTo));
                setIsLoading(false);
                setIsCalendarVisible(false);
            } catch (error) {
                notify('Could not read the collection contract from the network. Try later.');
                setIsCalendarVisible(false);
                setCalendarDateRange([new DateRangePick()]);
                setIsLoading(false);
            }
        })();
    }, [projectId, activeDropdownPeriod.period, unixDateFrom, unixDateTo]);

    /** Hides calendar and resets choosen period to initial state. */
    const cancelCalendarPeriod = () => {
        setCalendarDateRange([new DateRangePick()]);
        setActiveDropdownPeriod(chartPeriods[FIRST_ARRAY_INDEX]);
        setIsCalendarVisible(false);
    };

    /** Confirms calendar period and resets dropdown period to initial state. */
    const confirmCalendarPeriod = (dateRange: any) => {
        setCalendarDateRange([dateRange.selection]);
        setActiveDropdownPeriod(chartPeriods[FIRST_ARRAY_INDEX]);
    };

    return (
        <div className="dashboard__chart">
            {isLoading && <Loading />}
            <div className="dashboard__chart__header">
                <span className="dashboard__chart__label">Sales</span>
                <div className={`dashboard__chart__calendar ${calendarBlockClassName}`} >
                    <img
                        className="dashboard__chart__calendar__icon"
                        src={calendar}
                        alt="Calendar"
                        onClick={() => setIsCalendarVisible(!isCalendarVisible)}
                    />
                    <DateRange
                        editableDateInputs={true}
                        onChange={item => confirmCalendarPeriod(item)}
                        moveRangeOnFirstSelection={false}
                        ranges={calendarDateRange}
                        className={elementClassName}
                    />
                    <div className={`dashboard__chart__calendar__btn${elementClassName}`}>
                        <span className="choosen-days">{calendarPeriod + DAY} {daysDeclension}</span>
                        <button
                            className="cancel-btn"
                            onClick={cancelCalendarPeriod}
                        >
                            Cancel
                        </button>
                        <button
                            className="continue-btn"
                            onClick={() => setIsCalendarVisible(false)}
                        >
                            Continue
                        </button>
                    </div>
                </div>
                <Filters
                    activeDropdownPeriod={activeDropdownPeriod}
                    setActiveDropdownPeriod={setActiveDropdownPeriod}
                    setCalendarDateRange={setCalendarDateRange}
                    setIsCalendarVisible={setIsCalendarVisible}
                    chartPeriods={chartPeriods}
                    calendarPeriod={calendarPeriod}
                />
            </div>
            {!isLoading && !dashboardData.length &&
                <span className="dashboard__chart__empty">No sales in this period</span>
            }
            {!!dashboardData.length && <ResponsiveContainer width="100%">
                <AreaChart id="chart" data={dashboardData}>
                    <defs>
                        <linearGradient id="colorValue" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="0%" stopColor="#F8F5F5" />
                            <stop offset="39%" stopColor="#EDEDED" />
                            <stop offset="100%" stopColor="#F0F0F0" />
                        </linearGradient>
                    </defs>
                    <XAxis dy={AXIOS_PADDING} tickSize={CHART_STYLES_NUMBER} tick={{ stroke: '#9B9B9B', strokeWidth: CHART_STYLES_NUMBER }} dataKey="date" axisLine={false} />
                    <YAxis dx={-AXIOS_PADDING} tickSize={CHART_STYLES_NUMBER} tick={{ stroke: '#9B9B9B', strokeWidth: CHART_STYLES_NUMBER }} axisLine={false} type="number" domain={[CHART_STYLES_NUMBER, (dataMax: number) => dataMax * Y_AXIS_HEADROOM]} />
                    {/** @ts-ignore */}
                    <Tooltip content={<CustomTooltip />} />
                    <Area activeDot={{ r: DOTTER_RADIUS }} type="monotone" strokeWidth={LINE_STROKE} dataKey="value" stroke="#3DB9FF" fill="url(#colorValue)" />
                </AreaChart>
            </ResponsiveContainer>}
        </div>
    );
};
