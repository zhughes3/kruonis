import React, {FormEvent, useState} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {loginAttempt} from "../Http/Requests";
import isEmail from 'validator/lib/isEmail';
import {checkForPasswordLength} from "../Utils/PasswordChecks";

let email = '';
let password = '';
export const Login: React.FunctionComponent<IRouterProps> = (props) => {

	const [loginAttemptStatus, setLoginAttemptStatus] = useState<string>('');

	const handleLogin = async (e: FormEvent): Promise<void> => {
		e.preventDefault();

		setLoginAttemptStatus('');

		if (errorFound()) { return; }

		const result = await loginAttempt({email, password});
		console.log(result);

		resetRegisterData();
		// @ts-ignore
		document.getElementById("login-form")?.reset();

		props.history.push('/');
	}

	const resetRegisterData = (): void => {
		email = '';
		password = '';

		// @ts-ignore
		document.getElementById("login-form")?.reset();
	}

	const errorFound = (): boolean => {
		let error = false;

		if (!isEmail(email) || !checkForPasswordLength(password)) {
			error = true;
		}

		if (error) {
			setLoginAttemptStatus('Login failed, please check the data you entered')
		}

		return error;
	}

	return (
		<div className="login">
			<form id="login-form" onSubmit={handleLogin}>

				<div className="field has-text-centered">
					<div className="title">Login</div>
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
					<div className="control has-icons-left">
						<input className="input" type="password" placeholder="Password" onChange={ (e) => password = e.target.value} />
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

				{ loginAttemptStatus &&
				<div className="field has-text-danger">
					{loginAttemptStatus}
				</div>
				}

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
	)
}

