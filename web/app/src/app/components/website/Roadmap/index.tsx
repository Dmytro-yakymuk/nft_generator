import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import { RootState } from '@/app/store';
import { createMilestone, deleteMilestone, getMilestones, setMilestone, setMilestones, updateMilestone } from '@/app/store/actions/projects';
import { CreateMilestoneFields, Milestone } from '@/collections';
import { notify } from '@/internal/notifications';

import './index.scss';

// TODO: should be reworked with children props.
const Roadmap: React.FC = () => {
    const dispatch = useDispatch();
    const milestones = useSelector((state: RootState) => state.collectionsReducer.mintsite.roadmap);
    const { id }: { id: string } = useParams();

    /** Indicates counter. */
    const COUNTER_VALUE: number = 1;
    /** Default milestones list length. */
    const DEFAULT_MILESTONES_LENGTH: number = 1;
    /** Default social networks list length. */
    const FIRST_MILESTONES_LIST_ELEMENT_INDEX: number = 0;

    /** Adds new milestones. */
    const addMilestones = async() => {
        const milestoneFields = new CreateMilestoneFields();
        try {
            if (milestones.length === DEFAULT_MILESTONES_LENGTH) {
                const firstMilestoneFields = new CreateMilestoneFields(
                    milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].title,
                    milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].description,
                    milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].date,
                );
                milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].id ?
                    await dispatch(updateMilestone(milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX])) :
                    await dispatch(createMilestone(id, firstMilestoneFields));
                await dispatch(createMilestone(id, milestoneFields));
                await dispatch(getMilestones(id));

                return;
            }
            await Promise.all(
                milestones.map(
                    async(milestone: Milestone) => await dispatch(updateMilestone(milestone))
                )
            );
            await dispatch(createMilestone(id, milestoneFields));
            await dispatch(getMilestones(id));
            notify('Milestone saved!', 'success');
        } catch (error) {
            notify('Milestones not saved!');
        }
    };

    /** Saves milestones. */
    const saveMilestones = async() => {
        try {
            if (milestones.length === DEFAULT_MILESTONES_LENGTH) {
                const firstMilestoneFields = new CreateMilestoneFields(
                    milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].title,
                    milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].description,
                    milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].date,
                );
                milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX].id ?
                    await dispatch(updateMilestone(milestones[FIRST_MILESTONES_LIST_ELEMENT_INDEX])) :
                    await dispatch(createMilestone(id, firstMilestoneFields));
                await dispatch(getMilestones(id));

                return;
            }
            await Promise.all(
                milestones.map(
                    async(milestone: Milestone) => await dispatch(updateMilestone(milestone))
                )
            );
            await dispatch(getMilestones(id));
            notify('Milestones saved!', 'success');
        } catch (error) {
            notify('Milestones not saved!');
        }
    };

    /** Deletes milestone. */
    const deleteCurrentMilestone = (milestoneId: string) => {
        dispatch(deleteMilestone(id, milestoneId));
        const filteredMilestones = milestones.filter((milestone: Milestone) => milestone.id !== milestoneId);
        dispatch(setMilestones(filteredMilestones));
    };

    /** Changes milestone. */
    const changeMilestone = (id: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        let copiedMilestone = milestones.find((milestone: Milestone) => milestone.id === id);
        const index = milestones.findIndex((milestone: Milestone) => milestone.id === id);

        if (!copiedMilestone) {
            return;
        }

        copiedMilestone = { ...copiedMilestone, [name]: value };
        dispatch(setMilestone(index, copiedMilestone));
    };

    useEffect(() => {
        dispatch(getMilestones(id));
    }, []);

    return <section className="website">
        <div className="website__wrapper">
            <div className="website__roadmap-create">
                {milestones.map((milestone: Milestone, index: number) =>
                    <div
                        className="website__milestone"
                        key={milestone.id}
                    >
                        <h2 className="website__milestone__title">
                            Milestone {COUNTER_VALUE + index}
                        </h2>
                        <a
                            onClick={() => deleteCurrentMilestone(milestone.id)}
                            className="website__milestone__delete"
                            id={`btn-${index}`}
                        >
                            Delete
                        </a>
                        <label
                            htmlFor="title"
                            className="website__milestone__label-label"
                        >
                            Title
                        </label>
                        <input
                            className="website__milestone__label"
                            value={milestone.title}
                            onChange={changeMilestone(milestone.id)}
                            type="text"
                            id="title"
                            name="title"
                        />
                        <label
                            htmlFor="date"
                            className="website__milestone__date-label"
                        >
                            Date
                        </label>
                        <input
                            className="website__milestone__date"
                            value={milestone.date}
                            onChange={changeMilestone(milestone.id)}
                            type="text"
                            id="date"
                            name="date"
                        />
                        <label
                            htmlFor="description"
                            className="website__milestone__description-label"
                        >
                            Description (optional)
                        </label>
                        <input
                            className="website__milestone__description"
                            value={milestone.description}
                            onChange={changeMilestone(milestone.id)}
                            id="description"
                            name="description"
                            type="text"
                        />
                    </div>
                )}
                <div className="website__roadmap-create__buttons">
                    <button
                        onClick={addMilestones}
                        className="website__roadmap-create__new"
                    >
                        New Milestone
                    </button>
                    <button
                        onClick={saveMilestones}
                        className="website__roadmap-create__save"
                    >
                        Save
                    </button>
                </div>
            </div>
            <div className="website__custom-roadmap">
                <div className="website__custom-roadmap__wrapper">
                    <h2 className="website__custom-roadmap__title">
                        Roadmap
                    </h2>
                    <ul className="website__custom-roadmap__list">
                        {milestones.map((milestone: Milestone) =>
                            <li
                                className="website__custom-roadmap__item"
                                key={milestone.id}
                            >
                                <span className="website__custom-roadmap__item__date">
                                    {milestone.date}
                                </span>
                                <span className="website__custom-roadmap__item__title">
                                    {milestone.title}
                                </span>
                                <span className="website__custom-roadmap__item__decsription">
                                    {milestone.description}
                                </span>
                            </li>
                        )}
                    </ul>
                </div>
            </div>
        </div>
    </section>;
};

export default Roadmap;
