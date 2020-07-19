import React, {FunctionComponent, useEffect, useState} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {AbsoluteCenter} from "../Components/AbsoluteCenter";
import {deleteTimelineGroup, getUser, updateTimelineGroup} from "../Http/Requests";
import {IFullOrder} from "../Interfaces/IFullOrder";
import {IGroup} from "../Interfaces/IGroup";
import {Link} from "react-router-dom";

import edit from '../Assets/Icons/edit-2.svg';
import eye from "../Assets/Icons/eye.svg";
import eyeOff from "../Assets/Icons/eye-off.svg";
import trash from "../Assets/Icons/trash.svg";
import {LogoutButton} from "../Components/LogoutButton";
import {EditGroupName} from "../Components/EditGroupName";

export const Dashboard: FunctionComponent<IRouterProps> = (props) => {

	const [user, setUser] = useState<IFullOrder>();
	const [selectedGroup, setSelectedGroup] = useState<IGroup>();
	const [openEditGroupName, setOpenEditGroupName] = useState<boolean>(false);

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

					{ user && !user.groups &&
						<div className="columns is-gapless table-border-bottom table-item space-between">
							<div>Looks like you don't have any timelines yet. <Link to={'/'}>Let's go make some!</Link></div>
						</div>
					}

					{ user && user.groups &&
						user.groups.map( (group: IGroup) =>
							<div key={group.id} className="columns is-gapless table-border-bottom table-item space-between">

								<Link to={`timeline/${group.id}`}>{group.title}</Link>

								<div>
									<img className="cursor-pointer mr2" src={edit} alt="Private" title="Edit name" onClick={ () => { setSelectedGroup(group); setOpenEditGroupName(!openEditGroupName) } } />
									{ group.private ?
										<img className="cursor-pointer" src={eyeOff} alt="Private" title="Timeline is currently private" onClick={() => togglePrivatePublic(group, false)}/>
										:
										<img className="cursor-pointer" src={eye} alt="Public" title="Timeline is currently public" onClick={() => togglePrivatePublic(group, true)}/>
									}
									<img className="ml2 cursor-pointer" src={trash} alt="Remove" title="Remove timeline" onClick={ () => deleteGroup(group.id) } />
								</div>
							</div>
						)
					}
				</div>

				<EditGroupName open={openEditGroupName} group={selectedGroup} toggleModal={() => setOpenEditGroupName(!openEditGroupName)} updateTitle={ () => fetchMyData() }/>

				<LogoutButton className="mt4" />
			</div>
		</AbsoluteCenter>
	);
}
