import React from 'react';

interface ICenterProps {
    className?: string;
}

export const Center: React.FunctionComponent<ICenterProps> = (props) => {

    return (
        <div className={`${props.className}`} style={{display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
            <div>
                {props.children}
            </div>
        </div>
    );
}