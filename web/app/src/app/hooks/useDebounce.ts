import { useEffect, useState } from 'react';

/** Hook usedebounce allows to debounce any fast changing value over delay.*/
export const useDebounce = (value: string, delay: number) => {
    const [debouncedValue, setDebouncedValue] = useState(value);

    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedValue(value);
        }, delay);

        return () => {
            clearTimeout(timer);
        };
    }, [value]);

    return debouncedValue;
};
