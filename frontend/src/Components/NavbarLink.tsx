import React from 'react';
import {Link} from "react-router-dom";

interface INavbarLinkProps {
	name: string;
	to: string;
	icon?: string;
}

export const NavbarLink: React.FunctionComponent<INavbarLinkProps> = (props) => {

	return (
		<Link className="navbar-icon" style={{display: 'inline-flex', justifyContent: 'center'}} to={props.to}>
			{ props.icon && <img src={props.icon} alt={'Link to ' + props.name} /> }
			<span className="ml05 mr2" >{props.name}</span>
		</Link>
	);
}
