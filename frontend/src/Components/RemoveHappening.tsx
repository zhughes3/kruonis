import trash from "../Assets/Icons/trash.svg";
import React from "react";
import {deleteHappening} from "../Http/Requests";

interface IRemoveHappeningProps {
	// The id of the happening we want to delete.
	id: string;
}

export const RemoveHappening: React.FunctionComponent<IRemoveHappeningProps> = (props) => {

	const click = async () => {
		const result = await deleteHappening(props.id);
		console.log(result);
	}

	return <img src={trash} alt="Remove event" onClick={click} />;
}
