import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";

import { Home } from './Containers/Home';
import {Login} from "./Containers/Login";
import {Register} from "./Containers/Register";
import { Timeline } from "./Containers/Timeline";
import {NoMatch} from "./Containers/NoMatch";

export const Router = () => {
    return (
        <BrowserRouter>
            <div>
                <Switch>
                    <Route path="/" exact component={Home} />
                    <Route path="/login" exact component={Login} />
                    <Route path="/register" exact component={Register} />
                    <Route path="/timeline/:groupId" exact component={Timeline} />
                    <Route component={NoMatch} />
                </Switch>
            </div>
        </BrowserRouter>
    );
}
