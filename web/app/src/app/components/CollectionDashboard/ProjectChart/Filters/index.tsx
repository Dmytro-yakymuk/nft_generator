import { useState } from 'react';

import { ChartPeriod } from '@/dashboard';
import { DateRangePick } from '@components/CollectionDashboard/ProjectChart';

import './index.scss';

type FiltersProps = {
    activeDropdownPeriod: ChartPeriod;
    setActiveDropdownPeriod: (activeDropdownPeriod: ChartPeriod) => void;
    setIsCalendarVisible: (isCalendarVisible: boolean) => void;
    setCalendarDateRange: (dateRange: any | []) => void;
    chartPeriods: ChartPeriod[];
    calendarPeriod: number;
};

export const Filters: React.FC<FiltersProps> = ({ activeDropdownPeriod, setActiveDropdownPeriod, chartPeriods, calendarPeriod, setCalendarDateRange, setIsCalendarVisible }) => {
    const [isDropDownActive, setIsDropdownActive] = useState<boolean>(false);
    const dropdownClass: string = isDropDownActive ? '-active' : '';

    /** Selects chart period, closes calendar and sets choosen calendar period to initial state.
     * @param {ChartPeriod} period - date segment.
     */
    const handleChoosePeriod = (period: ChartPeriod): void => {
        setIsDropdownActive(false);
        setActiveDropdownPeriod(period);
        setCalendarDateRange([new DateRangePick()]);
        setIsCalendarVisible(false);
    };

    /** Method to change active filter class name.
     * @param {ChartPeriod} period - date segment.
     */
    const changeActiveDropdownItemClass = (period: ChartPeriod) => period.label === activeDropdownPeriod.label ? '-active': '';

    return (
        <div className="chart-filters">
            <div className="chart-filters__calendar"></div>
            <div className="chart-filters__dropdown">
                <span
                    className="chart-filters__dropdown__active-item"
                    onClick={() => setIsDropdownActive(!isDropDownActive)}
                >
                    {calendarPeriod ? 'Select Period' : activeDropdownPeriod.label}
                </span>
                <div className={`chart-filters__dropdown__list${dropdownClass}`}>
                    {
                        chartPeriods.map(
                            (period: ChartPeriod, index: number) =>
                                <span
                                    className={`chart-filters__dropdown__list__item${changeActiveDropdownItemClass(period)}`}
                                    key={index}
                                    onClick={() => handleChoosePeriod(period)}
                                >
                                    {period.label}
                                </span>
                        )
                    }
                </div>
            </div>
        </div>
    );
};
