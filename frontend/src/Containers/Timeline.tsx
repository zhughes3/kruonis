import React, {useEffect, useState} from 'react';
import { IRouterProps } from '../Interfaces/IRouterProps';
import { IHappening, IHappeningCreate } from '../Interfaces/IHappening';
import { Happening } from '../Components/Happening';
import { useParams } from 'react-router-dom';
import { createHappening, createImage, deleteHappening, getTimelineGroup, updateHappening } from "../Http/Requests";
import {Center} from "../Components/Center";
import { LoadingCard } from '../Components/LoadingCard';
import moment from "moment";
import {ErrorCard} from "../Components/ErrorCard";
import { IGroup } from '../Interfaces/IGroup';
import {ITimeline} from "../Interfaces/ITimeline";
import {HappeningModal} from "../Components/HappeningModal";
import edit from "../Assets/Icons/edit-2.svg";
import {TimelineTitles} from "../Components/TimelineTitles";

// This happening is displayed if there are no happenings (events) on the timeline yet.
const emptyTimelineHappening = {
    id: '1',
    timeline_id: '1',
    title: 'No events yet',
    timestamp: moment().format(),
    description: 'Click the "+" button (bottom left) to start adding items to your timeline.',
    content: "This is your new timeline. Right now, it's empty. But you can start to add items to it by clicking the '+' button.",
    created_at: 'created_at',
    updated_at: 'updated_at',
};

let events: any = [];
let selectedTimeline: any;

export const Timeline: React.FunctionComponent<IRouterProps> = (props) => {

    let { groupId } = useParams();

    const [open, setOpen] = useState<boolean>(false);
    const [selectedHappening, setSelectedHappening] = useState<IHappening | undefined>();
    const [loading, setLoading] = useState<boolean>(true);
    const [timelineGroup, setTimelineGroup] = useState<IGroup>();
    const [fetchTimelineError, setFetchTimelineError] = useState<string>('');
    const [openEditModal, setOpenEditModal] = useState<boolean>(false);

    useEffect( () => {
        if (!groupId) { return setFetchTimelineError("We can't find the timeline you are looking for!"); }
        fetchTimeLine(groupId);

        return undefined;
        // eslint-disable-next-line
    }, [groupId]);

    const fetchTimeLine = async (id: string): Promise<void> => {

        setLoading(true);

        const result = await getTimelineGroup(id).catch( async (e: Error) => {
            // This error may be coused because we are trying to view a private timeline while not logged in.
            // TODO implement proper behaviour when someone is trying to checkout a private timeline that does not belong to them.
            console.log(e);
        });

        if (result && result.id) {
            // The events are sorted and set in the setEvens function below.
            setEvents(result);

            setTimelineGroup(result);
        }

        // If a response has no .id, it's probably a 404 (no official 404).
        // Send user to home screen if timeline does not exist.
        // TODO implement message for user when navigating to timeline that doesn't exist.
        if (!result) {
            props.history.push('/');
        }

        setLoading(false);
    };

    const createNewHappening = async (happening: IHappeningCreate, timelineId: string): Promise<void> => {

        // TODO Add error handling on no (timelineId?) id.
        const result: IHappening | void = await createHappening(timelineId, happening).catch( (e: Error) => console.log(e) );

        if (!result) { return; }

        // If the happening contains an image, upload it.
        if (happening.image) {
            const image_url: { Url: string } = await createImage(result.id, happening.image);
            result.image_url = image_url.Url;
        }

        const timelineGroupCopy = Object.assign(timelineGroup, {});
        let timeIn = 0;

        const selectedTimeline = timelineGroupCopy.timelines.filter( (timeline, index) => {
            if (timeline.id !== timelineId) return false;
            timeIn = index;
            return true;
        })[0];

        // If a timeline has no events array we create one. This situation can occur when we delete all of the existing events off of an existing timeline with events.
        if (!selectedTimeline.events) { selectedTimeline.events = []; }

        selectedTimeline.events.push(result);

        timelineGroupCopy.timelines[timeIn].events = selectedTimeline.events;

        setEventsAndTimelineGroup(timelineGroupCopy);

    };

    const updateAHappening = async (happening: IHappening): Promise<void> => {

        // TODO Add error handling on no happening.
        const result: IHappening | void = await updateHappening(happening.id, happening).catch( (e: Error) => console.log(e) );

        if (!result) { return; }

        if (happening.image) {
            const image_url: { Url: string } = await createImage(result.id, happening.image);
            result.image_url = image_url.Url;
        }

        const cp = Object.assign({}, timelineGroup);

        cp?.timelines.forEach( timeline => {
            timeline.events.forEach( (event, index) => {
               if (event.id === result.id) {
                   timeline.events[index] = result;
               }
           });
        });

        setSelectedHappening(result);

        setEventsAndTimelineGroup(cp);
    };

    const removeHappening = async (id: string) => {

        // TODO Add error handling on delete failure.
        const result = await deleteHappening(id).catch( (e: Error) => console.log(e) );

        if (!result) { return; }

        const cp = Object.assign({}, timelineGroup);

        cp?.timelines.forEach( timeline => {

            if (!timeline.events) { return; }

            timeline.events.forEach( (event, index) => {
                if (event.id === id) {
                    timeline.events.splice(index, 1);
                }
            });
        });

        setSelectedHappening(undefined);

        setEventsAndTimelineGroup(cp);
    };

    const setEventsAndTimelineGroup = (newTimelineGroup: IGroup) => {
        setEvents(newTimelineGroup);

        setTimelineGroup(undefined);
        setTimelineGroup(timelineGroup => newTimelineGroup);
    };

    // Combine the events of multiple timelines into a single array. Then sort that array based on date.
    // This function is ALWAYS called after a post, put or deleter to make sure the events are displayed in the correct order.
    const setEvents = (group: IGroup) => {
        events = [];

        group.timelines.forEach( timeline => {

            if (!timeline.events) { return; }

            events = [...events, ...timeline.events];
        });

        events = sortHappenings(events);
    };

    const sortHappenings = (happenings: IHappening[]): IHappening[] => {
        return happenings.sort( (a, b) => {
            if (a.timestamp < b.timestamp) { return -1; }
            if ( a.timestamp > b.timestamp) { return 1; }
            return 0;
        });
    };

    // If the happening a user has clicked on is not active, display it. If it is, hide it.
    const toggleHappening = (happening: IHappening) => {
        if (selectedHappening?.id !== happening.id) {
            return setSelectedHappening(happening)
        }

        return setSelectedHappening(undefined);
    };

    const toggleModal = (): void => {
        setOpen(!open);
    };

    if (loading ) {
        return (
            <Center>
                <LoadingCard text="Loading your timeline..." />
            </Center>
        )
    }

    if (fetchTimelineError || !timelineGroup) {
        return (
            <Center>
                <ErrorCard text="Something went wrong when we tried to fetch your timeline." />
            </Center>
        )
    }

    return (
        <div>

            { timelineGroup?.timelines.map( (timeline: ITimeline, index: number) => {
                return (
                    <div key={index} className="add-happening" style={{marginBottom: index * 70, zIndex: 2, backgroundColor: 'white', padding: 5, borderRadius: 5}}>
                        <div className="fab has-background-link" onClick={() => {
                            toggleModal();
                            selectedTimeline = timeline.id
                        }}> + </div>
                        <div className="ml1">{timeline.title}</div>
                    </div>
                )
            })}

            {/*<AddHappening open={open} toggleModal={toggleModal} createHappening={ (newHappening: IHappeningCreate) => { createNewHappening(newHappening, selectedTimeline) }} />*/}
            {/* Modal used for creating new happenings. */}
            <HappeningModal  buttonText="Update" onSubmit={ (newHappening: IHappeningCreate) => { createNewHappening(newHappening, selectedTimeline) } } title="Create event"  open={open}  toggleModal={toggleModal} />

            <div className="timeline-position">

                {!events.length &&
                    <div className="notification is-link animated fadeIn fast" style={{maxWidth: 500, marginLeft: -283, marginBottom: 50}}>
                        <div><b>Welcome to your new timeline</b></div>
                        <div>To start adding events, click the "+" button on the bottom left of your screen.</div>
                    </div>
                }

                <TimelineTitles timelineGroup={timelineGroup} />

                {/* The timeline events. */}
                <div className="steps is-vertical is-centered is-small animated fadeInUp timeline-space">

                    {events.length ?
                        events.map((happening: IHappening, index: number) => {
                            return <Happening className="pb-70 cursor-pointer" key={index} left={happening.timeline_id === timelineGroup?.timelines[0].id} happening={happening} selectHappening={ (happening: IHappening) => toggleHappening(happening) }  openEditHappening={openEditModal} setOpenEditHappening={ (open) => setOpenEditModal(open) }  deleteHappening={ (id: string) => removeHappening(id) } />
                        })
                        :
                        <Happening className="pb-70 cursor-pointer" key={1} happening={emptyTimelineHappening} selectHappening={ (happening: IHappening) => toggleHappening(happening) } left openEditHappening={openEditModal} setOpenEditHappening={ (open) => setOpenEditModal(open) } deleteHappening={ (id: string) => removeHappening(id) } />
                    }

                </div>
            </div>

            <div className="happening-description">
                {selectedHappening &&
                    <div className="animated fadeInRight fast">
                        <div className="happening-info-title">{selectedHappening?.title}</div>
                        { selectedHappening?.image_url && <img className="mt2" style={{maxHeight: 400}} src={selectedHappening?.image_url} alt="Happening" /> }
                        <div className="mt-10 happening-info-content">{selectedHappening?.content}</div>
                    </div>
                }
            </div>

            {/* Modal used for updating happenings. */}
            <HappeningModal  buttonText="Update" onSubmit={ (result: IHappening) => {updateAHappening(result) }} title="Update event"  open={openEditModal}  toggleModal={ () => setOpenEditModal(!openEditModal) } happening={selectedHappening} />

        </div>
    );
};
