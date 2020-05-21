import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";

import { Home } from './Containers/Home';
import {Login} from "./Containers/Login";
import {Register} from "./Containers/Register";
import { Timeline } from "./Containers/Timeline";
import {NoMatch} from "./Containers/NoMatch";
import {Dashboard} from "./Containers/Dashboard";
import {PrivateRoute} from "./Components/PrivateRoute";

export const Router = () => {

    return (
        <BrowserRouter>
            <div>
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
}
