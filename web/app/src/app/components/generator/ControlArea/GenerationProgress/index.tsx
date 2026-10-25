import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Cell, Pie, PieChart } from 'recharts';
import { useDispatch, useSelector } from 'react-redux';

import { RootState } from '@/app/store';
import { abortGenerateNTTs, getGenerationProgress, setGenerationProgress } from '@/app/store/actions/projects';
import { notify } from '@/internal/notifications';

import './index.scss';

/** Exposes special type for */
class ChartValue {
    constructor(
        public value: number
    ) { }
}

/** Contains colors for loading progress chart */
class ColorsConfig {
    public static emptyPartColor: string = '#E6E6E6';
    public static loadingPartColor: string = '#000';
}

export const GenerationProgress: React.FC<{setGenerationPopUpVisibility: Dispatch<SetStateAction<boolean>>}> = ({ setGenerationPopUpVisibility }) => {
    const DEFAULT_VALUE = 0;
    const FIRST_ELEMENT_INDEX = 1;
    /** Equals to width and height of chart container in px */
    const CHART_AREA_DIMESSION = 148;
    /** Chart radius in px */
    const CHART_RADIUS = 74;
    /** Chart offset from chart container */
    const CHART_OFFSET = 70;
    /** Endpoint request interval */
    const INTERVAL = 500;
    const { id }: { id: string } = useParams();
    const dispatch = useDispatch();
    const nftCount = useSelector((state: RootState) => state.collectionsReducer.nftGenerationCount);
    const project = useSelector((state: RootState) => state.collectionsReducer.project);
    const [progressRest, setProgressRest] = useState(project.nftCount);

    /** Stops click event propagation for popup body */
    const handleStopPropagation = (e: React.MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();
    };

    useEffect(() => {
        const timer = setInterval(() => {
            dispatch(getGenerationProgress(id));
        }, INTERVAL);

        return () => {
            clearInterval(timer);
            dispatch(setGenerationProgress(DEFAULT_VALUE));
        };
    }, []);

    const loadingStatus = () => {
        if (nftCount < project.nftCount) {
            return 'Loading...';
        }

        return 'Done';
    };

    /** Cancels the generation: the images made so far stay, the collection can be generated again */
    const cancelGeneration = async() => {
        try {
            await abortGenerateNTTs(id);
            notify('Generation has been cancelled', 'success');
            setGenerationPopUpVisibility(false);
        } catch (error: any) {
            notify(`Could not cancel the generation: ${error.message}`);
        }
    };

    /** Sets calculated state for each request */
    useEffect(() => {
        setProgressRest(project.nftCount - nftCount);
    }, [nftCount]);

    return (
        <div className="generation-progress" onClick={() => setGenerationPopUpVisibility(false)}>
            <div className="generation-progress__wrapper" onClick={handleStopPropagation}>
                <h4 className="generation-progress__heading">
                    Generating Collection
                </h4>
                <span className="generation-progress__status">{loadingStatus()}</span>
                <PieChart width={CHART_AREA_DIMESSION} height={CHART_AREA_DIMESSION}>
                    <Pie
                        data={[new ChartValue(progressRest), new ChartValue(nftCount)]}
                        cx={CHART_OFFSET}
                        cy={CHART_OFFSET}
                        innerRadius={CHART_OFFSET}
                        outerRadius={CHART_RADIUS}
                        paddingAngle={DEFAULT_VALUE}
                        dataKey="value"
                    >
                        <Cell key={`cell-${DEFAULT_VALUE}`} fill={ColorsConfig.emptyPartColor} />
                        <Cell key={`cell-${FIRST_ELEMENT_INDEX}`} fill={ColorsConfig.loadingPartColor} />
                    </Pie>
                </PieChart>
                <div className="generation-progress__counter-area">
                    <span>Images</span>
                    <span>{`${nftCount}/${project.nftCount}`}</span>
                </div>
                {nftCount < project.nftCount &&
                    <button className="generation-progress__cancel" onClick={cancelGeneration}>
                        Cancel
                    </button>
                }
            </div>
        </div>
    );
};
