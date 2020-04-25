import React from 'react';
import moment from 'moment';

interface IHappeningDateProps {
    timestamp: string;
    className?: string;
    // Determines if this date rotated for a left side display.
    left?: boolean;
}

export const HappeningDate: React.FunctionComponent<IHappeningDateProps> = (props) => {

    return (
        <div className={`happening-date happening-minus-margin ${props.className && props.className} ${props.left ? 'rotate-for-left' : ''}`}>
            <div className={`happening-content ml2 ${props.left ? 'rotate-text-for-left' : ''}`}>
                <div className="month-day-text">{moment(props.timestamp).format('MMM Do')}</div>
                <div className="year-text">{moment(props.timestamp).format('YYYY')}</div>
            </div>
        </div>
    );
}
