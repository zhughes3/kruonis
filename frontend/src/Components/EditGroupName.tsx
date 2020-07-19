import React, {FormEvent, useEffect, useState} from 'react';
import {updateTimelineGroup} from "../Http/Requests";
import {IGroup} from "../Interfaces/IGroup";
import {ErrorCard} from "./ErrorCard";

interface IEditGroupNameProps {
	open: boolean;
	group: IGroup | undefined;
	toggleModal: () => void;
	updateTitle: () => void;
}

export const EditGroupName: React.FunctionComponent<IEditGroupNameProps> = (props) => {

	const [title, setTitle] = useState<string | undefined>(props.group ? props.group.title : '');

	useEffect(() => {
		setTitle(props.group ? props.group.title : '');
	}, [props.group])

	const handleSubmit = async (e: FormEvent) => {
		e.preventDefault();

		if (!props.group) { return; }

		const group = JSON.parse(JSON.stringify(props.group));
		group.title = title;

		await updateTimelineGroup(group);

		props.updateTitle();

		closeModal();
	}

	const closeModal = () => {
		props.toggleModal();
	}

	if (!props.group || !props.group.title) {
		return (
			<div className={`modal ${props.open && 'is-active'}`}>
				<div className="modal-background" onClick={() => closeModal()}/>
				<div className="modal-content">
					<ErrorCard text={'Error, no group selected'} />
				</div>
				<button className="modal-close is-large" aria-label="close" onClick={() => closeModal()} />
			</div>
		)
	}

	return (
		<div className={`modal ${props.open && 'is-active'}`}>
			<div className="modal-background" onClick={() => closeModal()}/>
			<div className="modal-content">
				<div className="box">

					<form onSubmit={handleSubmit} id="happening-form">
						<div className="field mt2">
							<label className="label">Name</label>
							<div className="control">
								<input className="input" type="text" placeholder="Event name" value={title} onChange={ (e) => setTitle( e.target.value ) } />
							</div>
						</div>

						<div className="field mt2">
							<button className="button is-success" type="submit">Save changes</button>
						</div>
					</form>
				</div>
			</div>
			<button className="modal-close is-large" aria-label="close" onClick={() => closeModal()} />
		</div>
	);
};
