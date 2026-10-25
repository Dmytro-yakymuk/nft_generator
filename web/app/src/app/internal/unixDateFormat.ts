const DATE_DIVIDER: number = 1000;
/** Converts date to unix format. */
export const convertDateToUnix = (date: Date) => Math.floor(date.getTime() / DATE_DIVIDER);
