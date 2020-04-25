import React, {useEffect, useState} from 'react';
import { IRouterProps } from '../Interfaces/IRouterProps';
import { IHappening, IHappeningCreate } from '../Interfaces/IHappening';
import { Happening } from '../Components/Happening';
import { AddHappening } from '../Components/AddHappening';
import { useParams } from 'react-router-dom';
import {ITimeline} from "../Interfaces/ITimeline";
import {createHappening, getTimeline} from "../Http/Requests";
import {Center} from "../Components/Center";
import { LoadingCard } from '../Components/LoadingCard';
import moment from "moment";
import {ErrorCard} from "../Components/ErrorCard";

// This happening is displayed if there are no happenings (events) on the timeline yet.
const emptyTimelineHappening = {
    id: '1',
    event_id: '1',
    timeline_id: '1',
    title: 'first event',
    timestamp: moment().format(),
    description: 'Click the "Add" button to start adding items to your timeline.',
    content: 'This is your new timeline. You can start to add items to it by clicking the "+" button.',
    created_at: 'created_at',
    updated_at: 'updated_at',
};

export const Timeline: React.FunctionComponent<IRouterProps> = (props) => {

    let { id } = useParams();

    const [open, setOpen] = useState<boolean>(false);
    const [selectedHappening, setSelectedHappening] = useState<IHappening>();
    const [loading, setLoading] = useState<boolean>(true);
    const [timeline, setTimeline] = useState<ITimeline>();
    const [fetchTimelineError, setFetchTimelineError] = useState<string>('');

    useEffect( () => {
        if (!id) { return setFetchTimelineError("We can't find the timeline you are looking for!"); }
        fetchTimeLine(id);
    }, [id]);

    const fetchTimeLine = async (id: string): Promise<void> => {

        setLoading(true);

        const result = await getTimeline(id).catch( (e: Error) => {
            console.log(e);
            setFetchTimelineError('Woops, something went wrong when we tried to fetch your timeline!');
            return setLoading(false);
        });

        if (result) {
            result.events = sortHappenings(result.events);
            setTimeline(result);
        }

        setLoading(false);
    };

    const createNewHappening = async (happening: IHappeningCreate): Promise<void> => {

        // TODO Add error handling on no id.
        if (!id) { return; }

        const result: IHappening | void = await createHappening(id, happening).catch( (e: Error) => console.log(e) );

        if (!result) { return; }

        const timelineCopy = Object.assign(timeline, {});
        timelineCopy.events.push(result);
        timelineCopy.events = sortHappenings(timelineCopy.events);

        setTimeline(undefined);
        setTimeline( timeline => timelineCopy);

    };

    const sortHappenings = (happenings: IHappening[]): IHappening[] => {
        // @ts-ignore
        return happenings.sort( (a, b) => a.timestamp > b.timestamp)
    };

    const toggleModal = (): void => {
        setOpen(!open);
    };

    if (loading) {
        return (
            <Center>
                <LoadingCard text="Loading your timeline..." />
            </Center>
        )
    }

    if (fetchTimelineError) {
        return (
            <Center>
                <ErrorCard text="Something went wrong when we tried to fetch your timeline." />
            </Center>
        )
    }

    return (
        <div>

            <button onClick={() => toggleModal()}>click</button>

            <AddHappening open={open} toggleModal={toggleModal} createHappening={ (newHappening: IHappeningCreate) => { createNewHappening(newHappening) }} />

            <div className="timeline-position">
                <div className="steps is-vertical is-centered is-small animated fadeInUp">

                    {timeline?.events.length ?
                        timeline?.events.map((happening: IHappening) => {
                            return <Happening className="pb-70 cursor-pointer" key={happening.event_id} happening={happening} selectHappening={ (happening: IHappening) => setSelectedHappening(happening) }/>
                        })
                        :
                        <Happening className="pb-70 cursor-pointer" key={1} happening={emptyTimelineHappening} selectHappening={ (happening: IHappening) => setSelectedHappening(happening) }/>
                    }

                </div>
            </div>

            <div className="happening-description">
                {selectedHappening &&
                    <div className="animated fadeInRight fast">
                        <div className="happening-info-title">{selectedHappening?.title}</div>
                        <div className="mt-10">{selectedHappening?.content}</div>
                    </div>
                }
            </div>

        </div>
    );
};
