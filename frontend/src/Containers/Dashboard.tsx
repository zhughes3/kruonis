import React, {FunctionComponent, useEffect} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {AbsoluteCenter} from "../Components/AbsoluteCenter";
import trash from "../Assets/Icons/trash.svg";

export const Dashboard: FunctionComponent<IRouterProps> = (props) => {

	useEffect(() => {
		// TODO fetch all timelines for a user.
	});

	return (
		<AbsoluteCenter>
			<div className="has-text-left">
				<div className="field">
					<div className="title">Timelines</div>
				</div>

				<div className="mt-36">
					<div className="columns is-gapless table-border-bottom table-item space-between">
						<div>WW 1 - WW 2</div>
						<img src={trash} alt="Remove event" onClick={ () => {} } />
					</div>
					<div className="columns is-gapless table-border-bottom table-item space-between">
						<div>The Netherlands - USA</div>
						<img src={trash} alt="Remove event" onClick={ () => {} } />
					</div>
					<div className="columns is-gapless table-item space-between">
						<div>Drake - Eminem</div>
						<img src={trash} alt="Remove event" onClick={ () => {} } />
					</div>
				</div>
			</div>
		</AbsoluteCenter>
	);
}
