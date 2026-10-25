const DAY: number = 24;
const HOUR_IN_SECONDS: number = 60;
const DATE_DIVIDER: number = 1000;
/** Calculates the difference between two dates.
 * @param {Date} startDate - Start date.
 * @param {Date} endDate - End date.
 */
export const calcDateDifference = (startDate: Date, endDate: Date): number => {
    const date1: Date = new Date(startDate);
    const date2: Date = new Date(endDate);
    const timeInMilisec: number = date2.getTime() - date1.getTime();
    const daysBetweenDates: number = Math.ceil(timeInMilisec / (DATE_DIVIDER * HOUR_IN_SECONDS * HOUR_IN_SECONDS * DAY));

    return daysBetweenDates;
};
