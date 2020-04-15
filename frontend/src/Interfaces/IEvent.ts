export interface IEvent {
    id: string;
    timeline_id: string;
    title: string;
    timestamp: string;
    description: string;
    content: string;
    created_at: string;
    updated_at: string;
}

export interface IEventCreate {
    timeline_id: string;
    title: string;
    timestamp: string;
    description: string;
    content: string;
}