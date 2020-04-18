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
        timestamp: '2010-06-13T18:25:43.511Z',
        description: 'This happend first',
        content: 'This is the content shown when the event is clicked',
        created_at: 'created_at',
        updated_at: 'updated_at',
    },
    {
        id: '2',
        timeline_id: '1',
        title: 'second event',
        timestamp: '2010-07-25T18:25:43.511Z',
        description: 'This happend second uawowohdnwa u ajwdnwaoddnwaodjn awuodjwanmdowadn wadowandkljwanmdoiwa donawwdnmwadnmmwadnw adwoddnwadidklwamdwa dwaujwalnd wawand wadjkwandjklawndoaw djwaanddjklwand wajwanddlaw nodwanldwa',
        content: 'This is the second item shown when the event is clicked',
        created_at: 'created_at',
        updated_at: 'updated_at',
    },
    {
        id: '3',
        timeline_id: '1',
        title: 'third event',
        timestamp: '2010-07-25T18:25:43.511Z',
        description: 'This happend third uawowohdnwa u ajwdnwaoddnwaodjn awuodjwanmdowadn wadowandkljwanmdoiwa donawwdnmwadnmmwadnw adwoddnwadidklwamdwa dwaujwalnd wawand wadjkwandjklawndoaw djwaanddjklwand wajwanddlaw nodwanldwa',
        content: 'This is the second item shown when the event is clicked',
        created_at: 'created_at',
        updated_at: 'updated_at',
    },
    {
        id: '4',
        timeline_id: '1',
        title: 'fourth event',
        timestamp: '2010-07-25T18:25:43.511Z',
        description: 'Another short story is written over here now.',
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

            <div className="absolute-center">
                <div className="steps is-vertical is-centered is-small animated fadeInUp">

                    {happening.map((happening: IHappening) => {
                        return <Happening className="pb-70" key={happening.id} happening={happening} />
                    })}

                </div>
            </div>

        </div>
    );
}