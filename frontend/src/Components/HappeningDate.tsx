import React from 'react';
import moment from 'moment';

interface IHappeningDateProps {
    timestamp: string;
    className?: string;
}

export const HappeningDate: React.FunctionComponent<IHappeningDateProps> = (props) => {

    return (
        <div className={`happening-date happening-minus-margin ${props.className && props.className}`}>
            <div className="happening-content ml2">
                <div className="month-day-text">{moment(props.timestamp).format('MMM Do')}</div>
                <div className="year-text">{moment(props.timestamp).format('YYYY')}</div>
            </div>
        </div>
    );
}