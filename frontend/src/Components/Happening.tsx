import React from 'react';
import { IHappening } from '../Interfaces/IHappening';
import { HappeningDate } from './HappeningDate';

interface IHappeningProps {
    happening: IHappening;
}

export const Happening: React.FunctionComponent<IHappeningProps> = (props) => {

    return (
        <div style={{display: 'flex', width: 300, alignItems: 'center', justifyContent: 'center'}}>
            
            <div>
                <div>{props.happening.title}</div>
                <div>{props.happening.description}</div>
            </div>
            
            <div>
                <HappeningDate timestamp={props.happening.timestamp} />
            </div>
        </div>
    );
}