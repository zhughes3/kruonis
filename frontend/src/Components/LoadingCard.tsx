import React from 'react';
import {Spinner} from "./Spinner";

interface ILoadingCard {
	text: string;
}

export const LoadingCard: React.FunctionComponent<ILoadingCard> = (props) => {

	return (
		<div className="loading-center">

			<div>

				<div style={{marginLeft: 60, marginBottom: 16}}>
					<Spinner />
				</div>

				<div>{props.text}</div>

			</div>

		</div>
	);
}
