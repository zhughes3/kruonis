import React from 'react';
import { IHappening } from '../Interfaces/IHappening';
import { HappeningDate } from './HappeningDate';

interface IHappeningProps {
    happening: IHappening;
}

export const Happening: React.FunctionComponent<IHappeningProps> = (props) => {

    return (
        <li className="steps-segment">
            <span className="steps-marker"></span>
            <div className="steps-content">
                <HappeningDate timestamp={props.happening.timestamp} />
                <p className="is-size-4">{props.happening.title}</p>
                <p>{props.happening.description}</p>
            </div>
        </li>
    );
}