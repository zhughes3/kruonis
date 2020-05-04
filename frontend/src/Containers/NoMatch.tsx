import React from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {Center} from "../Components/Center";

export const NoMatch: React.FunctionComponent<IRouterProps> = (props) => {

	return (
		<Center>
			<div className="absolute-center-parent">
				<div className="absolute-center-child has-text-centered">

					<div className="field">
						<div className="title">404! Woops, looks like you made a wrong turn somewhere!</div>
					</div>

				</div>
			</div>
		</Center>
	)
}

