import React, { useState } from 'react';
import { IRouterProps } from '../Interfaces/IRouterProps';
import { IHappening } from '../Interfaces/IHappening';
import { Happening } from '../Components/Happening';
import { AddHappening } from '../Components/AddHappening';

const happening: IHappening[] = [
    {
        id: '1',
        timeline_id: '1',
        title: 'first event',
        timestamp: '12/04/2020',
        description: 'This happend first',
        content: 'This is the content shown when the event is clicked',
        created_at: 'created_at',
        updated_at: 'updated_at',
    },
    {
        id: '2',
        timeline_id: '1',
        title: 'second event',
        timestamp: '15/04/2020',
        description: 'This happend second',
        content: 'This is the second item shown when the event is clicked',
        created_at: 'created_at',
        updated_at: 'updated_at',
    }
]

export const Timeline: React.FunctionComponent<IRouterProps> = (props) => {

    const [open, setOpen] = useState<boolean>(false);

    const toggleModal = (): void => {
        setOpen(!open);
    }

    return (
        <div>
            The timeline page

            <button onClick={() => toggleModal()}>click</button>

            <AddHappening open={open} toggleModal={toggleModal} />

            <ul className="steps is-vertical is-centered is-small" style={{minHeight: 140 * happening.length}}>

                {happening.map((happening: IHappening) => {
                    return <Happening key={happening.id} happening={happening} />
                })}

            </ul>

        </div>
    );
}