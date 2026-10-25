import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Link, useParams } from 'react-router-dom';

import { RootState } from '@/app/store';
import { getProject } from '@store/actions/collections';
import { notify } from '@/internal/notifications';

import { ProjectValues } from './ProjectValues';
import { ProjectChart } from './ProjectChart';
import { ProjectInfo } from './ProjectInfo';
import { Activity } from './Activity';

import './index.scss';

const ZERO_ADDRESS = '0x0000000000000000000000000000000000000000';

/** Collection dashboard: statistics, sales chart and activity built from the events of the collection contract. */
export const ProjectDashboard = () => {
    const dispatch = useDispatch();
    const { id }: { id: string } = useParams();
    const project = useSelector((state: RootState) => state.collectionsReducer.project);
    const [isProjectLoaded, setIsProjectLoaded] = useState<boolean>(false);

    useEffect(() => {
        (async() => {
            try {
                await dispatch(getProject(id));
                setIsProjectLoaded(true);
            } catch (error) {
                notify('Something wrong. Try later.');
            }
        })();
    }, [id]);

    if (!isProjectLoaded || project.id !== id) {
        return <div className="dashboard" />;
    }

    if (project.contractAddress.toLowerCase() === ZERO_ADDRESS) {
        return (
            <div className="dashboard">
                <div className="dashboard__empty">
                    <span className="dashboard__empty__title">{project.name}</span>
                    <p className="dashboard__empty__text">
                        The collection contract is not deployed yet. The dashboard shows mints and sales
                        from the contract events once it is deployed.
                    </p>
                    <Link className="dashboard__empty__link" to={`/create-contract/${id}`}>Deploy the contract</Link>
                </div>
            </div>
        );
    }

    return (
        <div className="dashboard">
            <div className="dashboard__wrapper">
                <div className="dashboard__aside">
                    <ProjectInfo projectId={id} />
                </div>
                <div className="dashboard__content">
                    <ProjectChart projectId={id} />
                    <ProjectValues />
                    <Activity projectId={id} />
                </div>
            </div>
        </div>
    );
};
