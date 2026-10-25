import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import image from '@static/images/dashboard/image.svg';
import person from '@static/images/dashboard/person.svg';
import sales from '@static/images/dashboard/sales.svg';
import projectLogo from '@static/images/dashboard/project-logo.svg';

import { notify } from '@/internal/notifications';
import { RootState } from '@/app/store';
import { getCollectionStatistics } from '@/app/store/actions/dashboard';

import './index.scss';

/** Project info entity. */
class ProjectStatistic {
    public constructor(
        public logo: string,
        public value: number,
        public title: string,
    ) {}
}

export const ProjectInfo: React.FC<{ projectId: string }> = ({ projectId }) => {
    const dispatch = useDispatch();
    const collectionStatistics = useSelector((state: RootState) => state.dashboardReducer.collectionStatistics);
    const project = useSelector((state: RootState) => state.collectionsReducer.project);

    /** Project collection statistics. */
    const projectInfo: ProjectStatistic[] = [
        new ProjectStatistic(image, collectionStatistics.items, 'Items'),
        new ProjectStatistic(person, collectionStatistics.owners, 'Owners'),
        new ProjectStatistic(sales, collectionStatistics.sales, 'Sales'),
    ];

    useEffect(() => {
        (async() => {
            try {
                await dispatch(getCollectionStatistics(projectId));
            } catch (error) {
                notify('Could not read the collection contract from the network. Try later.');
            }
        })();
    }, [projectId]);

    return (
        <div className="dashboard__info">
            <div className="dashboard__info__wrapper">
                <img className="dashboard__info__logo" src={projectLogo} alt="Project logo" />
                <span className="dashboard__info__name">{project.name}</span>
            </div>
            <div className="dashboard__info__list">
                {projectInfo.map((item: ProjectStatistic, index: number) =>
                    <div className="dashboard__info__list__item" key={index}>
                        <img className="dashboard__info__list__item__logo" src={item.logo} alt="Image" />
                        <div className="dashboard__info__list__item__wrapper">
                            <span className="dashboard__info__list__item__value">{item.value}</span>
                            <span className="dashboard__info__list__item__title">{item.title}</span>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};
