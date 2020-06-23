import React, {useContext, useEffect} from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";

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

    // When we open the website, check if a user is logged in. If they are, get the user data and put in the global store.
    useEffect(() => {
        checkIfLoggedIn().then( async (result: boolean) => {
            if(!result) { return; }
            const user = await getUser();
            userStore.setUser(user);
        }).catch( (e: Error) => {
            console.log('Error: ', e)
        });
    }, []);

    return (
        <BrowserRouter>
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
        </BrowserRouter>
    );
});
