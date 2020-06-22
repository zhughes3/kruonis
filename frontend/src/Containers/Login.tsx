import React, {FormEvent, useContext, useState} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {getUser, loginAttempt} from "../Http/Requests";
import isEmail from 'validator/lib/isEmail';
import {checkForPasswordLength} from "../Utils/PasswordChecks";
import {observer} from "mobx-react";
import {UserStoreContext} from "../Store/UserStore";

let email = '';
let password = '';
export const Login: React.FunctionComponent<IRouterProps> = observer( (props) => {

	const userStore = useContext(UserStoreContext)

	const [loginAttemptStatus, setLoginAttemptStatus] = useState<string>('');

	const handleLogin = async (e: FormEvent): Promise<void> => {
		e.preventDefault();

		setLoginAttemptStatus('');

		if (errorFound()) { return; }

		// Returns a 401 on faulty login.
		const result = await loginAttempt({email, password}).catch(e => console.log(e));

		const user = await getUser();

		userStore.setUser(user);

		resetRegisterData();
		// @ts-ignore
		document.getElementById("login-form")?.reset();

		if (props.location.state && props.location.state.from) {
			return props.history.push(props.location.state.from);
		}

		return props.history.push('/dashboard');
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
							<span className="has-text-right has-text-grey-dark cursor-pointer" onClick={ () => props.history.push('reset-password') }>
								Forgot password?
							</span>
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
						<span className="cursor-pointer" onClick={ () => props.history.push('register') }>
							<span className="icon is-small is-left">
								<i className="fas fa-user" />
							</span>
							<span className="ml1 has-text-grey-dark">Create new account</span>
						</span>
					</p>
				</div>
			</form>
		</div>
	)
});

