import {ITimeline} from "./ITimeline";

export interface IGroup {
	id: string;
	title: string;
	timelines: ITimeline[];
	created_at: string;
	updated_at: string;
}
