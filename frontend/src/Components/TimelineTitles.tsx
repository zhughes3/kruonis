import React, {useState} from 'react';
import edit from "../Assets/Icons/edit-2.svg";
import {ITimeline} from "../Interfaces/ITimeline";
import {IGroup} from "../Interfaces/IGroup";
import {EditGroupName} from "./EditGroupName";

interface ITimelineTitlesProps {
	timelineGroup: IGroup;
}

export const TimelineTitles: React.FunctionComponent<ITimelineTitlesProps> = (props) => {

	const [open, setOpen] = useState<boolean>(false);
	const [group, setGroup] = useState<IGroup>(props.timelineGroup);

	const updateTitles = (group: IGroup | undefined): void => {
		if (!group) { return; }


	}

	return (
		<div className="timeline_titles_div" style={{marginLeft: -315}}>
			<div className="columns is-gapless space-between">
				<div className="timeline_title animated fadeInUp faster">{props.timelineGroup?.title}</div>
				<img src={edit} className="edit_timeline_names ml1 cursor-pointer" alt="Edit event" onClick={ () => setOpen(!open) } />
			</div>

			{/* The timeline titles. */}
			<div className="columns is-gapless">
				{ props.timelineGroup?.timelines.map( (timeLine: ITimeline, index: number) => {
					if (index === 0) {
						return <div key={index} className="timeline_name animated fadeInUp faster" style={{width: 320}}>{timeLine?.title}</div>
					}
					return <div key={index} className="timeline_name animated fadeInUp faster" style={{marginLeft: 160}}>{timeLine?.title}</div>
				})}
			</div>

			<EditGroupName open={open} group={props.timelineGroup} toggleModal={() => setOpen(!open)} updateTitle={ (group: IGroup | undefined) => updateTitles(group) } />
		</div>
	);
};
