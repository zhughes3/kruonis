import {IGroup} from "./IGroup";
import {IUser} from "./IUser";

export interface IFullOrder {
	user:   IUser;
	groups: IGroup[];
}
