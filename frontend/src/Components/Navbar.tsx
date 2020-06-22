import React, {useContext} from 'react';
import {UserStoreContext} from "../Store/UserStore";
import {observer} from "mobx-react";
import {NavbarLink} from "./NavbarLink";
import home from './../Assets/Icons/home.svg';
import user from './../Assets/Icons/user.svg';

export const Navbar: React.FunctionComponent = observer( () => {

	const userStore = useContext(UserStoreContext)

	return (
		<div className='navigation-bar'>
			<NavbarLink name="Home" to={'/'} icon={home} />
			<NavbarLink name={ userStore.user ? userStore.user.user.email : 'Account' } to={'/dashboard'} icon={user} />
		</div>
	);
});
