import React, {FormEvent} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {Center} from "../Components/Center";
import {signUpAttempt} from "../Http/Requests";

let email = '';
let password = '';
let repeatPassword = '';

export const Register: React.FunctionComponent<IRouterProps> = (props) => {

	const handleRegister = async (e: FormEvent): Promise<void> => {
		e.preventDefault();

		if (!checkIfPasswordsMatch(password, repeatPassword)) { return; }

		const result = await signUpAttempt({email, password});

		if(result && result.response) {
			resetRegisterData();
			props.history.push('login');
		}
	}

	const checkIfPasswordsMatch = (password: string, repeatPassword: string) => {
		return password === repeatPassword;
	}

	const resetRegisterData = (): void => {
		email = '';
		password = '';
		repeatPassword = '';

		// @ts-ignore
		document.getElementById("login-form")?.reset();
	}

	return (
		<Center>
			<div className="login">
				<form id="login-form" onSubmit={handleRegister}>

					<div className="field has-text-centered">
						<div className="title">Register</div>
					</div>

					<div className="field mt3">
						<p className="control has-icons-left has-icons-right">
							<input className="input" type="email" placeholder="Email" onChange={ (e) => email = e.target.value} />
							<span className="icon is-small is-left">
								<i className="fas fa-envelope"/>
							</span>
						</p>
					</div>

					<div className="field mt2">
						<p className="control has-icons-left">
							<input className="input" type="password" placeholder="Password" onChange={ (e) => password = e.target.value} />
							<span className="icon is-small is-left">
								<i className="fas fa-lock"/>
							</span>
						</p>
					</div>

					<div className="field mt2">
						<p className="control has-icons-left">
							<input className="input" type="password" placeholder="Repeat password" onChange={ (e) => repeatPassword = e.target.value} />
							<span className="icon is-small is-left">
								<i className="fas fa-lock"/>
							</span>
						</p>
					</div>

					<div className="field mt3">
						<p className="control">
							<button className="button house-blue-button is-fullwidth">Login</button>
						</p>
					</div>

					<div className="field mt2">
						<p className="has-text-centered">
							<a onClick={ () => props.history.push('login') }>
								<span className="icon is-small is-left">
								  	<i className="fas fa-arrow-left" />
								</span>
								<span className="ml1 has-text-grey-dark">Back to login</span>
							</a>
						</p>
					</div>
				</form>
			</div>
		</Center>
	)
}

