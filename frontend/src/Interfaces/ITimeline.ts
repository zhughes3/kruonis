import { IEvent } from "./IEvent";

export interface ITimeline {
    id: string;
    group_id: string;
    title: string;
    tags: string[];
    events: IEvent[],
    created_at: string;
    updated_at: string;
}

export interface ITimelineCreate {
    title: string;
    tags: string[];
    group_id?: number | string;
}