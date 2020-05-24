import React, {FunctionComponent, useEffect, useState} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {AbsoluteCenter} from "../Components/AbsoluteCenter";
import trash from "../Assets/Icons/trash.svg";
import {deleteTimelineGroup, getUser} from "../Http/Requests";
import {IFullOrder} from "../Interfaces/IFullOrder";
import {IGroup} from "../Interfaces/IGroup";
import {Link} from "react-router-dom";

export const Dashboard: FunctionComponent<IRouterProps> = (props) => {

	const [user, setUser] = useState<IFullOrder>();

	useEffect(() => {
		fetchMyData();
	}, []);

	const fetchMyData = async () => {
		const result = await getUser();
		setUser(result);
	};

	const deleteGroup = async (id: string): Promise<void> => {
		await deleteTimelineGroup(id).catch( (e: Error) => console.log(e) );
		fetchMyData();
	};

	return (
		<AbsoluteCenter>
			<div className="has-text-left">
				<div className="field">
					<div className="title">Timelines</div>
				</div>

				<div className="mt-36">

					{ user && user.groups &&
						user.groups.map( (group: IGroup) =>
							<div className="columns is-gapless table-border-bottom table-item space-between">
								<Link to={`timeline/${group.id}`}>{group.title}</Link>
								<img src={trash} alt="Remove event" onClick={ () => deleteGroup(group.id) } />
							</div>
						)
					}
				</div>
			</div>
		</AbsoluteCenter>
	);
}
