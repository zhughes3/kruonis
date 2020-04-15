import React from 'react';

interface IAddHappeningProps {
    open: boolean;
    toggleModal: () => void;
}

export const AddHappening: React.FunctionComponent<IAddHappeningProps> = (props) => {

    return (
        <div className={`modal ${ props.open && 'is-active' }`}>
            <div className="modal-background" onClick={props.toggleModal} />
            <div className="modal-content">
                <div className="box">
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean efficitur sit amet massa fringilla egestas. Nullam condimentum luctus turpis.
                </div>
            </div>
            <button className="modal-close is-large" aria-label="close" onClick={props.toggleModal} />
        </div>
    );
}