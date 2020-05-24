import React from 'react';
import {Link} from "react-router-dom";

export const Navbar: React.FunctionComponent = () => {

	return (
		<div className='navigation-bar'>
			<Link className="mr2" to={'/'}>Home</Link>
			<Link to={'/dashboard'}>Account</Link>
		</div>
	);
}
