import React from 'react';
import { IHappening } from '../Interfaces/IHappening';
import { HappeningDate } from './HappeningDate';

interface IHappeningProps {
    happening: IHappening;
    className?: string;
}

export const Happening: React.FunctionComponent<IHappeningProps> = (props) => {

    return (
        <div id={props.happening.id} className={`steps-segment ${props.className && props.className}`}>
            
            <span className="steps-marker"></span>
            
            <div className="steps-content columns">
                
                <HappeningDate timestamp={props.happening.timestamp} />
                
                <div className="ml2 happening-info has-text-centered">
                    <div className="happening-info-title">{props.happening.title}</div>
                    <div className="happening-info-title">-</div>
                    <div className="happening-info-text">{props.happening.description}</div>
                </div>
            </div>
        </div>
    );
}