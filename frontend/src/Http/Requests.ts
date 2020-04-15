import { ITimeline, ITimelineCreate } from "../Interfaces/ITimeline";
import { httpGet, httpPost } from "./HttpSetup";

export const getTimeline = async (id: string): Promise<ITimeline> => {
    return httpGet('timelines/' + id);
}

export const createTimeline = async (body: ITimelineCreate): Promise<ITimeline> => {
    return httpPost('timelines', body);
}

export const createGroupedTimelines = async (timeline: ITimelineCreate, compareTo: ITimelineCreate): Promise<ITimeline[]> => {
    const firstTimeline: ITimeline = await httpPost('timelines', timeline);

    compareTo.group_id = firstTimeline.group_id;

    const secondTimeline: ITimeline = await httpPost('timelines', compareTo);

    return [firstTimeline, secondTimeline]
}