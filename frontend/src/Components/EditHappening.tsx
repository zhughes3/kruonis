import React, {useState} from "react";
import edit from "../Assets/Icons/edit-2.svg";
import {IHappening} from "../Interfaces/IHappening";

interface IEditHappeningProps {
	happening: IHappening;
}

export const EditHappening: React.FunctionComponent<IEditHappeningProps> = (props) => {

	const [isActive, setIsActive] = useState<boolean>(false);

	const click = async () => {
		// const result = await deleteHappening(props.id);
		// console.log(result);
	}

	return (
		<div>
			<img src={edit} className="ml1" alt="Edit event" onClick={click} />

			<div className="modal">
				<div className="modal-background" />
				<div className="modal-content">
					Update modal.
				</div>
				<button className="modal-close is-large" aria-label="close" />
			</div>
		</div>
	);
};
