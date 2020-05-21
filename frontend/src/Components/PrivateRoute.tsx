import {Redirect, Route} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {checkIfLoggedIn} from "../Http/Requests";

/**
 * @EXPLANATION
 * This component checks if a user is logged in. If they are, send the user to the page they want to go to.
 * If they are not, redirect that user back to the login page.
 *
 */
// @ts-ignore
export const PrivateRoute = ({ component, ...options }) => {

	const [loading, setLoading] = useState<boolean>(true);
	const [loggedIn, setLoggedIn] = useState<boolean>(false);

	useEffect(() => {

		checkIfLoggedIn().then( result => {

			// If an error occurs.
			if(!result) {
				setLoggedIn(false);
				setLoading(false);
				return;
			}

			// On success.
			setLoggedIn(true);
			setLoading(false);
		}).catch( (e: Error) => {
			setLoggedIn(false);

			setLoading(false);
			console.log('Error: ', e)
		});
	}, []);

	// Displayed while loading.
	if (!loggedIn && loading) {
		return <div />
	}

	// If login check came back false, send user to login. We also send the page of origin so we can redirect back to the private page after login.
	if (!loggedIn && !loading) {
		return <Redirect {...options} to={{pathname: "/login", state: { from: options.path }}} />;
	}

	// If login check came back true, allow the user to go to the appropriate route.
	return <Route {...options} component={component} />;
};
