import './index.scss';

/** First array index. */
const FIRST_INDEX = 0;

export const CustomTooltip: React.FC<{ active: any; payload: any; label: string }>= ({ active, payload, label }) => {
    if (active && payload && payload.length) {
        return (
            <div className="custom-tooltip">
                <p className="label">{`${payload[FIRST_INDEX].value} ETH`}</p>
            </div>
        );
    }

    return null;
};
