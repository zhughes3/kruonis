import {ITimeline, ITimelineCreate} from "../Interfaces/ITimeline";
import {httpDelete, httpGet, httpPost, httpPut} from "./HttpSetup";
import {IHappening, IHappeningCreate} from "../Interfaces/IHappening";
import {IGroup} from "../Interfaces/IGroup";
import {IUserCreate} from "../Interfaces/IUser";
import {IBoolResponse} from "../Interfaces/IBoolResponse";
import {IFullOrder} from "../Interfaces/IFullOrder";

export const getTimelineGroup = async (groupId: string): Promise<IGroup> => {
    return httpGet('groups/' + groupId);
};

export const createTimeline = async (body: ITimelineCreate): Promise<ITimeline> => {
    return await httpPost('timelines', body);
};

export const createGroupedTimelines = async (timeline: ITimelineCreate, compareTo: ITimelineCreate): Promise<ITimeline[]> => {
    // Create the first timeline.
    const firstTimeline: ITimeline = await httpPost('timelines', timeline);

    // The first timeline returned a group_id, add that to the second timeline.
    compareTo.group_id = firstTimeline.group_id;

    // Now that the group id has been set, create the second timeline.
    const secondTimeline: ITimeline = await httpPost('timelines', compareTo);

    // Finally, return both timelines.
    return [firstTimeline, secondTimeline]
};

export const updateTimelineGroup = async (group: IGroup): Promise<IGroup> => {
    return await httpPut('groups/' + group.id, group);
}

export const deleteTimelineGroup = async (timelineGroupId: string): Promise<void> => {
    return await httpDelete('groups/' + timelineGroupId);
}

export const createHappening = async (timelineId: string, body: IHappeningCreate): Promise<IHappening> => {
    return await httpPost('timelines/' + timelineId + '/events', body);
};

export const updateHappening = async (eventId: string, body: IHappeningCreate): Promise<IHappening> => {
    return await httpPut('events/' + eventId, body);
};

export const deleteHappening = async (happeningId: string): Promise<any> => {
    return await httpDelete('events/' + happeningId);
}

export const createImage = async (eventId: string, image: File): Promise<{Url: string}> => {
    return await httpPost('events/' + eventId + '/img', image, true);
}

// Returns true if the register attempt was successful.
export const signUpAttempt = async (body: IUserCreate): Promise<boolean> => {
    const result: IBoolResponse = await httpPost('users/signup', body)
    return result.response;
}

// Returns true if login was successful. Returns a 401 on faulty login.
export const loginAttempt = async (body: IUserCreate): Promise<any> => {
    return await httpPost('users/login', body);
}

export const logoutAttempt = async (): Promise<void> => {
    return await httpGet('users/logout');
}

// Returns true if login was successful.
export const checkIfLoggedIn = async (): Promise<boolean> => {
    const result: IBoolResponse = await httpGet('users/ping').catch( (e: Error) => console.log(e) );
    if (!result) { return false; }
    return result.response;
}

export const getUser = async (): Promise<IFullOrder> => {
    return await httpGet('users/me');
}
