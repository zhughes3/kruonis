import React from 'react';
import { IHappening } from '../Interfaces/IHappening';
import { HappeningDate } from './HappeningDate';

interface IHappeningProps {
    happening: IHappening;
    className?: string;
    // This boolean determines whether the happening info should be displayed on the left or right side of the line.
    left?: boolean;

    selectHappening: (happening: IHappening) => void;
}

export const Happening: React.FunctionComponent<IHappeningProps> = (props) => {

    console.log(props.left);

    return (
        <div id={props.happening.id} className={`steps-segment ${props.className && props.className}`} onClick={ () => props.selectHappening(props.happening) }>

            { props.left ?
                <div>
                    <span className="steps-marker"></span>

                    <div className="steps-content columns" style={{position: 'relative', right: 406}}>

                        <div className="ml2 happening-info has-text-centered">
                            <div className="happening-info-title">{props.happening.title}</div>
                            <div className="happening-info-title">-</div>
                            <div className="happening-info-text">{props.happening.description}</div>
                        </div>

                        <HappeningDate timestamp={props.happening.timestamp} left />
                    </div>
                </div>
                :
                <div>
                    <span className="steps-marker"></span>

                    <div className="steps-content columns">

                        <HappeningDate timestamp={props.happening.timestamp}/>

                        <div className="ml2 happening-info has-text-centered">
                            <div className="happening-info-title">{props.happening.title}</div>
                            <div className="happening-info-title">-</div>
                            <div className="happening-info-text">{props.happening.description}</div>
                        </div>
                    </div>
                </div>

            }

        </div>
    );
}
