import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";

import { Home } from './Containers/Home';
import { Timeline } from "./Containers/Timeline";

export const Router = () => {
    return (
        <BrowserRouter>
            <div>
                <Switch>
                    <Route path="/" exact component={Home} />
                    <Route path="/timeline/:groupId" exact component={Timeline} />
                </Switch>
            </div>
        </BrowserRouter>
    );
}
