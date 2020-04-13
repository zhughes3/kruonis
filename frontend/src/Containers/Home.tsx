import React, { FormEvent } from 'react';
import { IRouterProps } from '../Interfaces/IRouterProps';
import isAlphanumeric from 'validator/lib/isAlphanumeric';

import home_blue_stain from './../Assets/home_blue_stain.svg';

let timelineName: string = '';
let timelineCompareName: string = '';

export const Home: React.FunctionComponent<IRouterProps> = (props) => {

    const handleSubmit = async (e: FormEvent): Promise<void> => {
        e.preventDefault();

        // Check if we have a name, remove the spaces from name so we can check if it contains of legal characters.
        if (!timelineName || !legalName(timelineName) || (timelineCompareName && !legalName(timelineCompareName))) {
            return;
        }

        const timeline = await createTimeline(timelineName, timelineCompareName);

        props.history.push(`timeline/${timeline.id}`);
    }

    const legalName = (name: string): boolean => {
        return isAlphanumeric(name.split(" ").join(''))
    }

    const createTimeline = async (timelineName: string, timelineCompareName: string): Promise<any> => {

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

                    <button className="home-form-button mt-36 animated fadeInUp"><span style={{ marginLeft: '39%' }}>Let's start!</span></button>
                </form>
            </div>

            <div className="home-right-side">
                <div className="home-subheader animated fadeInUp faster">Create Interactive Timelines</div>
                <div className="home-right-text animated fadeInUp">Create interactive timelines quickly and easily. No registration needed!</div>
            </div>
        </div>
    );
}