/** Converts date string to US date format.
 * @param {string} date - Date string value.
 */
export const convertDate = (date: string) => {
    const currentDate = new Date(date);
    const dateOptions: any = { day: 'numeric', month: 'short', year: 'numeric' };

    return currentDate.toLocaleDateString('en-GB', dateOptions);
};

/** Converts date string to date and time, e.g. 11 Oct 2026, 13:53: several events often happen on one day.
 * @param {string} date - Date string value.
 */
export const convertDateTime = (date: string) => {
    const currentDate = new Date(date);
    const dateOptions: any = { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' };

    return currentDate.toLocaleString('en-GB', dateOptions);
};
