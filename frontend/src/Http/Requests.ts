import {ITimeline, ITimelineCreate} from "../Interfaces/ITimeline";
import {httpGet, httpPost} from "./HttpSetup";
import {IHappening, IHappeningCreate} from "../Interfaces/IHappening";
import {IGroup} from "../Interfaces/IGroup";

export const getTimeline = async (id: string): Promise<ITimeline> => {
    return httpGet('timelines/' + id);
};

export const getTimelineGroup = async (groupId: string): Promise<IGroup> => {
    return httpGet('groups/' + groupId);
};

export const createTimeline = async (body: ITimelineCreate): Promise<ITimeline> => {
    return await httpPost('timelines', body);
};

export const createHappening = async (timelineId: string, body: IHappeningCreate): Promise<IHappening> => {
    return await httpPost('timelines/' + timelineId + '/events', body);
};

export const createGroupedTimelines = async (timeline: ITimelineCreate, compareTo: ITimelineCreate): Promise<ITimeline[]> => {
    const firstTimeline: ITimeline = await httpPost('timelines', timeline);

    compareTo.group_id = firstTimeline.group_id;

    const secondTimeline: ITimeline = await httpPost('timelines', compareTo);

    return [firstTimeline, secondTimeline]
};
