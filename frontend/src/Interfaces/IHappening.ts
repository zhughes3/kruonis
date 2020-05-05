import {ReactText} from "react";

export interface IHappening {
    id: string;
    event_id: string;
    timeline_id: string;
    title: string;
    timestamp: string;
    description: string;
    content: string;
    created_at: string;
    updated_at: string;
}

export interface IHappeningCreate {
    id?: ReactText;
    event_id?: ReactText;
    timeline_id?: string;
    title: string;
    timestamp: string;
    description: string;
    content: string;
}
