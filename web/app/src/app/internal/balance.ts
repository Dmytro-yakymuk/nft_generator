/**
  * Rounds the value to the given number of decimal places.
  * @param {number} number - Value to be rounded.
  * @param {number} floatPartLength - Number of characters after the decimal point.
  */
export function toFixedNoRounding(number: number, floatPartLength: number): number {
    const DEFAULT_VALUE = 0;
    const INDEX_DIFFERENCE = 1;

    if (!number.toString().includes('.')) {
        return number;
    }

    return +number.toString().slice(DEFAULT_VALUE, number.toString().indexOf('.') + floatPartLength + INDEX_DIFFERENCE);
};
