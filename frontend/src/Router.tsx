import React, {useContext, useEffect} from "react";
import { Switch, Route, useLocation } from "react-router-dom";

import { Home } from './Containers/Home';
import { Login } from "./Containers/Login";
import { Register } from "./Containers/Register";
import { Timeline } from "./Containers/Timeline";
import { NoMatch } from "./Containers/NoMatch";
import { Dashboard } from "./Containers/Dashboard";
import { PrivateRoute } from "./Components/PrivateRoute";
import { Navbar } from "./Components/Navbar";
import {checkIfLoggedIn, getUser} from "./Http/Requests";
import {observer} from "mobx-react";
import {UserStoreContext} from "./Store/UserStore";

export const Router = observer( () => {

    const userStore = useContext(UserStoreContext)
    let location = useLocation();

    // When we open the website and on every page change, check if a user is logged in. If they are, get the user data and put in the global store. Otherwise remove the user data from the global store.
    useEffect(() => {
        checkIfLoggedIn().then( async (result: boolean) => {
            // If the user is not logged in, set the user in the global state to undefined.
            if(!result) {
                userStore.setUser(undefined);
                return;
            }

            // if the user is logged in, but no user is stored in the global state, store it.
            if (!userStore.user) {
                const user = await getUser();
                userStore.setUser(user);
            }
        // If an error occurs, we set the user in the global state to undefined, but we don't attempt a logout.
        }).catch( (e: Error) => {
            console.log('Error: ', e)
            userStore.setUser(undefined);
        });
    }, [location]);

    return (
        <div>

            <Navbar />

            <Switch>
                <Route path="/" exact component={Home} />
                <Route path="/login" component={Login} />
                <Route path="/register" component={Register} />
                <Route path="/timeline/:groupId" component={Timeline} />
                <PrivateRoute path="/dashboard" component={Dashboard} />
                <Route component={NoMatch} />
            </Switch>

        </div>
    );
});
