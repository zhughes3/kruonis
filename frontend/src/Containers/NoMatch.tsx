import React from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {AbsoluteCenter} from "../Components/AbsoluteCenter";

export const NoMatch: React.FunctionComponent<IRouterProps> = (props) => {

	return (
		<AbsoluteCenter>
			<div className="field">
				<div className="title">404! Woops, looks like you made a wrong turn somewhere!</div>
			</div>
		</AbsoluteCenter>
	)
}

