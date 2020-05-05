import React from 'react';

export const AbsoluteCenter: React.FunctionComponent = (props) => {

	return (
		<div className="absolute-center-parent">
			<div className="absolute-center-child has-text-centered">

				{props.children}

			</div>
		</div>
	)
}

