import {ReactText} from "react";

export interface IHappening {
    id: string;
    timeline_id: string;
    title: string;
    timestamp: string;
    description: string;
    content: string;
    created_at: string;
    updated_at: string;
    image_url?: string;
    // IHappening can have an image file during on update on the front end (when selecting a new image).
    image?: File;
}

export interface IHappeningCreate {
    id?: ReactText;
    // event_id?: ReactText;
    timeline_id?: string;
    title: string;
    timestamp: string;
    description: string;
    content: string;
    image?: File;
}
