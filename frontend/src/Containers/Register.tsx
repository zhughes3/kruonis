import React, {FormEvent, useState} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {signUpAttempt} from "../Http/Requests";
import isEmail from 'validator/lib/isEmail';
import { checkForEmptyPassword, checkForPasswordLength, checkIfPasswordsMatch, MIN_PASSWORD_LENGTH } from "../Utils/PasswordChecks";

let email = '';
let password = '';
let repeatPassword = '';

export const Register: React.FunctionComponent<IRouterProps> = (props) => {

	const [passwordMismatch, setPasswordMismatch] = useState<string>('');
	const [validEmail, setValidEmail] = useState<string>('');

	const handleRegister = async (e: FormEvent): Promise<void> => {
		e.preventDefault();

		resetErrors();

		if (errorFound()) { return; }

		const result = await signUpAttempt({email, password});

		if(result) {
			resetRegisterData();
			props.history.push('login');
		}
	}

	const errorFound = (): boolean => {
		let error = false;

		if (!checkIfPasswordsMatch(password, repeatPassword)) {
			setPasswordMismatch('Passwords do not match');
			error = true;
		}

		if (!checkForPasswordLength(password)) {
			setPasswordMismatch(`Password needs to be at least ${MIN_PASSWORD_LENGTH} characters long`);
			error = true;
		}

		if (!checkForEmptyPassword(password)) {
			setPasswordMismatch('Please enter a password');
			error = true;
		}

		if (!isEmail(email)) {
			setValidEmail('Please enter a valid email');
			error = true;
		}

		return error;
	}

	const resetRegisterData = (): void => {
		email = '';
		password = '';
		repeatPassword = '';

		// @ts-ignore
		document.getElementById("login-form")?.reset();
	}

	const resetErrors = (): void => {
		setPasswordMismatch('');
		setValidEmail('');
	}

	return (
		<div className="login">
			<form id="login-form" onSubmit={handleRegister}>

				<div className="field has-text-centered">
					<div className="title">Register</div>
				</div>

				<div className="field mt3">
					<p className="control has-icons-left has-icons-right">
						<input className={`input ${ validEmail && 'is-danger' }`} type="email" placeholder="Email" onChange={ (e) => email = e.target.value} />
						<span className="icon is-small is-left">
							<i className="fas fa-envelope"/>
						</span>
					</p>
				</div>

				{ validEmail &&
					<div className="field has-text-danger">
						{validEmail}
					</div>
				}

				<div className="field mt2">
					<p className="control has-icons-left">
						<input className={`input ${ passwordMismatch && 'is-danger' }`} type="password" placeholder="Password" onChange={ (e) => password = e.target.value} />
						<span className="icon is-small is-left">
							<i className="fas fa-lock"/>
						</span>
					</p>
				</div>

				<div className="field mt2">
					<p className="control has-icons-left">
						<input className={`input ${ passwordMismatch && 'is-danger' }`} type="password" placeholder="Repeat password" onChange={ (e) => repeatPassword = e.target.value} />
						<span className="icon is-small is-left">
							<i className="fas fa-lock"/>
						</span>
					</p>
				</div>

				{ passwordMismatch &&
					<div className="field has-text-danger">
						{passwordMismatch}
					</div>
				}

				<div className="field mt3">
					<p className="control">
						<button className="button house-blue-button is-fullwidth">Register</button>
					</p>
				</div>

				<div className="field mt2">
					<p className="has-text-centered">
						<span className="cursor-pointer" onClick={ () => props.history.push('login') }>
							<span className="icon is-small is-left">
								<i className="fas fa-arrow-left" />
							</span>
							<span className="ml1 has-text-grey-dark">Back to login</span>
						</span>
					</p>
				</div>
			</form>
		</div>
	)
}

