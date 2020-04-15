import React, { FormEvent, useState } from 'react';
import { IRouterProps } from '../Interfaces/IRouterProps';
import isAlphanumeric from 'validator/lib/isAlphanumeric';

import home_blue_stain from './../Assets/home_blue_stain.svg';
import { createTimeline, createGroupedTimelines } from '../Http/Requests';
import { ITimeline } from '../Interfaces/ITimeline';
import { Spinner } from '../Components/Spinner';

let timelineName: string = '';
let timelineCompareName: string = '';

export const Home: React.FunctionComponent<IRouterProps> = (props) => {

    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>('');

    const handleSubmit = async (e: FormEvent): Promise<void> => {
        e.preventDefault();

        // Check if we have a name, remove the spaces from name so we can check if it contains illegal characters.
        if (!timelineName || !legalName(timelineName) || (timelineCompareName && !legalName(timelineCompareName))) {
            return setError("Don't forget to name your timeline ( Letters and numbers only )!");
        }

        try {

            setLoading(true);

            let createResult: ITimeline | ITimeline[];

            if (!timelineCompareName) {
                createResult = await createTimeline({title: timelineName, tags: []});
                props.history.push(`timeline/${createResult.id}`);
            } else {
                createResult = await createGroupedTimelines({title: timelineName, tags: []}, {title: timelineCompareName, tags: []});
                props.history.push(`timeline/${createResult[0].group_id}`);
            }
        }
        catch(e) {
            console.log(e);
            setError('Woops, something went wrong!');
            setLoading(false);
        }
    }

    const legalName = (name: string): boolean => {
        return isAlphanumeric(name.split(" ").join(''))
    }

    return (
        <div>

            <img className="home_blue_stain" src={home_blue_stain} alt="" />

            <div className="home-left-side">
                <div className="home_intro_text color-white animated fadeInUp faster">Create your timeline</div>

                <form className="create-timeline-form" onSubmit={handleSubmit}>
                    <div className="home-form-text color-white animated fadeInUp faster">Your timeline name:</div>
                    <input className="home-form-input mt-17 animated fadeInUp faster" type="text" placeholder="e.g Eminem" onChange={(e) => timelineName = e.target.value} />

                    <div className="home-form-text mt-36 animated fadeInUp">( Optional ) Compare your timeline to:</div>
                    <input className="home-form-input mt-17 animated fadeInUp" type="text" placeholder="e.g Drake" onChange={(e) => timelineCompareName = e.target.value} />

                    <button className="home-form-button mt-36 animated fadeInUp"><span style={{ marginLeft: loading ? '28%' : '39%' }}>{loading ? "Creating timeline" : "Let's start!"} {loading && <Spinner />}</span></button>
                </form>

                { error && <div className="mt-10 pb-10 color-error font-size-18">{error}</div> }
            </div>

            <div className="home-right-side">
                <div className="home-subheader animated fadeInUp faster">Create Interactive Timelines</div>
                <div className="home-right-text animated fadeInUp">Create interactive timelines quickly and easily. No registration needed!</div>
            </div>
        </div>
    );
}