import React, {useContext} from 'react';
import {logoutAttempt} from "../Http/Requests";
import { useHistory } from "react-router-dom";
import {UserStoreContext} from "../Store/UserStore";

interface ILogoutButtonProps {
	className?: string;
}

export const LogoutButton: React.FunctionComponent<ILogoutButtonProps> = (props) => {

	let history = useHistory();
	const userStore = useContext(UserStoreContext)

	const attemptLogout = async (): Promise<void> => {
		await logoutAttempt().catch( () => {
			// Remove the user from the global storage.
			userStore.setUser(undefined);
			// Send the user to the home page.
			history.push("/");
		})
	};

	return (
		<button style={{background: '#2251CA'}} className={`button pt1 pb1 pl3 pr3 is-link ${props.className}`} onClick={() => attemptLogout()}>Logout</button>
	);
}
