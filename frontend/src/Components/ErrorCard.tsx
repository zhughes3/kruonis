import React from 'react';

interface IErrorCard {
	text: string;
}

export const ErrorCard: React.FunctionComponent<IErrorCard> = (props) => {

	return (
		<div className="error-center">

			<div>

				<div>{props.text}</div>

			</div>

		</div>
	);
}
