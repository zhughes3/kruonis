import React from 'react';

interface IHappeningDateProps {
    timestamp: string;
}

export const HappeningDate: React.FunctionComponent<IHappeningDateProps> = (props) => {

    return (
        <div>
            {props.timestamp};
        </div>
    );
}