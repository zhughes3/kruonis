import React, {FormEvent} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {Center} from "../Components/Center";
import {loginAttempt, signUpAttempt} from "../Http/Requests";

let email = '';
let password = '';
export const Login: React.FunctionComponent<IRouterProps> = (props) => {

	const handleLogin = async (e: FormEvent): Promise<void> => {
		e.preventDefault();

		const result = await loginAttempt({email, password});
		console.log(result);

		// @ts-ignore
		document.getElementById("login-form")?.reset();
	}

	const resetRegisterData = (): void => {
		email = '';
		password = '';

		// @ts-ignore
		document.getElementById("login-form")?.reset();
	}

	return (
		<Center>
			<div className="login">
				<form id="login-form" onSubmit={handleLogin}>

					<div className="field has-text-centered">
						<div className="title">Login</div>
					</div>

					<div className="field mt3">
						<p className="control has-icons-left has-icons-right">
							<input className="input" type="email" placeholder="Email"/>
							<span className="icon is-small is-left">
								<i className="fas fa-envelope"/>
							</span>
						</p>
					</div>

					<div className="field mt2">
						<div className="control has-icons-left">
							<input className="input" type="password" placeholder="Password"/>
							<span className="icon is-small is-left">
								<i className="fas fa-lock"/>
							</span>
							<div className="has-text-right mt1">
								<a className="has-text-right has-text-grey-dark" onClick={ () => props.history.push('reset-password') }>
									Forgot password?
								</a>
							</div>
						</div>
					</div>

					<div className="field mt3">
						<p className="control">
							<button className="button house-blue-button is-fullwidth">Login</button>
						</p>
					</div>

					<div className="field mt2">
						<p className="has-text-centered">
							<a onClick={ () => props.history.push('register') }>
								<span className="icon is-small is-left">
								  	<i className="fas fa-user" />
								</span>
								<span className="ml1 has-text-grey-dark">Create new account</span>
							</a>
						</p>
					</div>
				</form>
			</div>
		</Center>
	)
}

