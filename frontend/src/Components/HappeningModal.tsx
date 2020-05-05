import React, {FormEvent, useEffect, useState} from 'react';
import {IHappening, IHappeningCreate} from '../Interfaces/IHappening';
import moment from "moment";
import edit from "../Assets/Icons/edit-2.svg";

interface IHappeningModalProps {
	open: boolean;

	// Type of onSubmit: IHappeningCreate | IHappening.
	onSubmit: (newHappening: any) => void;

	toggleModal: () => void;
	title: string;
	buttonText: string;

	happening?: IHappening;
}

export const HappeningModal: React.FunctionComponent<IHappeningModalProps> = (props) => {

	const [errorMessage, setErrorMessage] = useState<string>('');

	const [title, setTitle] = useState<string>('');
	const [description, setDescription] = useState<string>('');
	const [content, setContent] = useState<string>('');
	const [timestamp, setTimestamp] = useState<string>('');

	useEffect(() => {
		if (props.happening) {
			setTitle(props.happening.title);
			setDescription(props.happening.description);
			setContent(props.happening.content);
			setTimestamp(props.happening.timestamp.split('T')[0])
		}
	}, [props.happening, props.open])

	const handleSubmit = (e: FormEvent) => {
		e.preventDefault();

		if (!title || !description || !timestamp) {
			return setErrorMessage('Please fill in atleast a title, date and short description')
		}

		setErrorMessage('');

		const happeningResult: IHappeningCreate = {
			id: props.happening ? props.happening.id : 0,
			// event_id is needed for updating.
			event_id: props.happening ? props.happening.event_id : 0,
			title,
			timestamp: moment.utc(timestamp).format(),
			description,
			content,
		};

		props.onSubmit(happeningResult);
		closeModal();
	};

	const closeModal = (): void => {
		// Clear the form and input data.
		setTitle('');
		setDescription('');
		setContent('');
		setTimestamp('');
		// @ts-ignore
		document.getElementById("happening-form")?.reset();

		setErrorMessage('');
		props.toggleModal();
	};

	return (
		<div className={`modal ${props.open && 'is-active'}`}>
			<div className="modal-background" onClick={() => closeModal()}/>
			<div className="modal-content">
				<div className="box">

					<form onSubmit={handleSubmit} id="happening-form">

						<div className="field">
							<div className="title">{props.title}</div>
						</div>

						<div className="field mt2">
							<label className="label">Name</label>
							<div className="control">
								<input className="input" type="text" placeholder="Event name" value={title}
									   onChange={(e) => setTitle(e.target.value) }/>
							</div>
						</div>

						<div className="field">
							<label className="label">Date</label>
							<div className="control">
								<input className="input" type="date" name="event-date-input" id="event-date" value={timestamp}
									   onChange={(e) => setTimestamp(e.target.value)} />
							</div>
						</div>

						<div className="field">
							<label className="label">Short description</label>
							<div className="control">
								<input className="input" type="text" placeholder="A short description" value={description}
									   onChange={(e) => setDescription(e.target.value) }/>
							</div>
						</div>

						<div className="field">
							<label className="label">More information</label>
							<div className="control">
								<textarea className="textarea"
										  value={content}
										  placeholder="This information is displayed when you click the event"
										  onChange={(e) => setContent(e.target.value) }/>
							</div>
						</div>

						{errorMessage && <div className="field color-error font-size-18">{errorMessage}</div>}

						<div className="field">
							<button className="button is-success">{props.buttonText}</button>
						</div>
					</form>
				</div>
			</div>
			<button className="modal-close is-large" aria-label="close" onClick={() => closeModal()}/>
		</div>
	);
};
