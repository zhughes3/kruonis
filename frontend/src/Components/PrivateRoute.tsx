import {Redirect, Route} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {checkIfLoggedIn} from "../Http/Requests";

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

	// If login check came back false.
	if (!loggedIn && !loading) {
		return <Redirect {...options} to={{pathname: "/login", state: { from: options.path }}} />;
	}

	// If login check came back true;
	return <Route {...options} component={component} />;
};
