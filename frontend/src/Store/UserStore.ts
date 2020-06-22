import { observable } from "mobx"
import {createContext} from "react";
import {IUser} from "../Interfaces/IUser";
import {IFullOrder} from "../Interfaces/IFullOrder";

class UserStore {
	@observable user: IFullOrder | undefined = undefined;

	setUser(user: IFullOrder | undefined): IFullOrder | undefined {
		this.user = user;
		return this.user;
	}
}

export const UserStoreContext = createContext(new UserStore());
