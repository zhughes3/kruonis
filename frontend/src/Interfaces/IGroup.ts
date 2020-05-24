import {ITimeline} from "./ITimeline";

export interface IGroup {
	id:         string;
	title:      string;
	timelines:  ITimeline[];
	created_at: Date;
	updated_at: Date;
	private:    boolean;
	user_id:    string;
	uuid:       string;
}
