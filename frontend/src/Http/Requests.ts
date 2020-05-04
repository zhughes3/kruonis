import {ITimeline, ITimelineCreate} from "../Interfaces/ITimeline";
import {httpDelete, httpGet, httpPost} from "./HttpSetup";
import {IHappening, IHappeningCreate} from "../Interfaces/IHappening";
import {IGroup} from "../Interfaces/IGroup";
import {IUserCreate} from "../Interfaces/IUser";

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

export const deleteHappening = async (happeningId: string): Promise<any> => {
    return await httpDelete('events/' + happeningId);
}

export const signUpAttempt = async (body: IUserCreate): Promise<{ response: boolean }> => {
    return await httpPost('users/signup', body)
}

// Returns a JWT token if successful.
export const loginAttempt = async (body: IUserCreate): Promise<{ token: string }> => {
    return await httpPost('users/login', body)
}
