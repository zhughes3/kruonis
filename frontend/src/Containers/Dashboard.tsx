import React, {FunctionComponent, useEffect, useState} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {AbsoluteCenter} from "../Components/AbsoluteCenter";
import {deleteTimelineGroup, getUser, updateTimelineGroup} from "../Http/Requests";
import {IFullOrder} from "../Interfaces/IFullOrder";
import {IGroup} from "../Interfaces/IGroup";
import {Link} from "react-router-dom";

import eye from "../Assets/Icons/eye.svg";
import eyeOff from "../Assets/Icons/eye-off.svg";
import trash from "../Assets/Icons/trash.svg";

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

	const togglePrivatePublic = async (group: IGroup, isPrivate: boolean): Promise<void> => {
		group.private = isPrivate;
		const result = await updateTimelineGroup(group);

		let tempUser = Object.assign({}, user);

		tempUser.groups = tempUser?.groups.map( g => {
			if (g.id === result.id) {
				g = result;
			}

			return g;
		});

		setUser(tempUser);
	}

	return (
		<AbsoluteCenter>
			<div className="has-text-left">
				<div className="field">
					<div className="title">Timelines</div>
				</div>

				<div className="mt-36">

					{ user && user.groups &&
						user.groups.map( (group: IGroup) =>
							<div key={group.id} className="columns is-gapless table-border-bottom table-item space-between">
								<Link to={`timeline/${group.id}`}>{group.title}</Link>
								<div>
									{ group.private ?
										<img src={eyeOff} alt="Private" onClick={() => togglePrivatePublic(group, false)}/>
										:
										<img src={eye} alt="Public" onClick={() => togglePrivatePublic(group, true)}/>
									}
									<img className="ml2" src={trash} alt="Remove" onClick={ () => deleteGroup(group.id) } />
								</div>
							</div>
						)
					}
				</div>
			</div>
		</AbsoluteCenter>
	);
}
