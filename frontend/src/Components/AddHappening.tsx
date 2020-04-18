import React, { FormEvent } from 'react';
import { IHappeningCreate } from '../Interfaces/IHappening';

interface IAddHappeningProps {
    open: boolean;
    createHappening: (newHappening: IHappeningCreate) => void;
    toggleModal: () => void;
}

let title: string = '';
let description: string = '';
let content: string = '';
let timestamp: string = '';

export const AddHappening: React.FunctionComponent<IAddHappeningProps> = (props) => {

    const handleSubmit = (e: FormEvent) => {
        e.preventDefault();

        const newHappening: IHappeningCreate = {
            title,
            timestamp,
            description,
            content
        }

        props.createHappening(newHappening);
    }

    return (
        <div className={`modal ${props.open && 'is-active'}`}>
            <div className="modal-background" onClick={props.toggleModal} />
            <div className="modal-content">
                <div className="box">

                    <form onSubmit={handleSubmit}>
                    
                        <div className="field">
                            <div className="title">Add event</div>
                        </div>

                        <div className="field mt2">
                            <label className="label">Name</label>
                            <div className="control">
                                <input className="input" type="text" placeholder="Event name" onChange={ (e) => title = e.target.value} />
                            </div>
                        </div>

                        <div className="field">
                            <label className="label">Date</label>
                            <div className="control">
                                <input className="input" type="date" name="event-date-input" id="event-date" onChange={ (e) => timestamp = e.target.value} />
                            </div>
                        </div>

                        <div className="field">
                            <label className="label">Short description</label>
                            <div className="control">
                                <input className="input" type="text" placeholder="A short description" onChange={ (e) => description = e.target.value} />
                            </div>
                        </div>

                        <div className="field">
                            <label className="label">More information</label>
                            <div className="control">
                                <textarea className="textarea" placeholder="This information is displayed when you click the event" onChange={ (e) => content = e.target.value} />
                            </div>
                        </div>

                        <div className="field">
                            <button className="button is-success">Add event</button>
                        </div>
                    </form>
                </div>
            </div>
            <button className="modal-close is-large" aria-label="close" onClick={props.toggleModal} />
        </div>
    );
}