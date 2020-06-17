import React from 'react';
import { IHappening } from '../Interfaces/IHappening';
import { HappeningDate } from './HappeningDate';

import trash from './../Assets/Icons/trash.svg';
import edit from './../Assets/Icons/edit-2.svg';

interface IHappeningProps {
    happening: IHappening;
    className?: string;
    // This boolean determines whether the happening info should be displayed on the left or right side of the line.
    left?: boolean;

    selectHappening: (happening: IHappening) => void;
    openEditHappening: boolean;
    setOpenEditHappening: (openEditHappening: boolean) => void;

    deleteHappening: (id: string) => void;
}

export const Happening: React.FunctionComponent<IHappeningProps> = (props) => {

    return (
        <div id={props.happening.id} className={`steps-segment ${props.className && props.className}`} onClick={ () => props.selectHappening(props.happening) }>

            { props.left ?
                <div className="happening">

                    {/* The icons for deleting and updating a happening. These become visible on hover. */}
                    <div className="happening-update-delete-left">
                        <img src={trash} alt="Remove event" onClick={ () => props.deleteHappening(props.happening.id) } />
                        <img src={edit} className="ml1" alt="Edit event" onClick={ () => props.setOpenEditHappening(!props.openEditHappening) } />
                    </div>

                    <span className="steps-marker" />

                    <div className="steps-content columns" style={{position: 'relative', right: 406}}>

                        <div className="ml2 happening-info has-text-centered">
                            <div className="happening-info-title">{props.happening.title}</div>
                            <div className="happening-info-title">-</div>
                            <div className="happening-info-text">{props.happening.description}</div>
                        </div>

                        <HappeningDate timestamp={props.happening.timestamp} left/>
                    </div>
                </div>
                :
                <div className="happening">

                    {/* The icons for deleting and updating a happening. These become visible on hover. */}
                    <div className="happening-update-delete-right">
                        <img src={trash} alt="Remove event" onClick={ () => props.deleteHappening(props.happening.id) } />
                        <img src={edit} className="ml1" alt="Edit event" onClick={ () => props.setOpenEditHappening(!props.openEditHappening) } />
                    </div>

                    <span className="steps-marker" />

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
