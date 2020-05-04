import React, {FormEvent} from 'react';
import {IRouterProps} from "../Interfaces/IRouterProps";
import {Center} from "../Components/Center";

export const Login: React.FunctionComponent<IRouterProps> = (props) => {

	const handleLogin = (e: FormEvent): void => {
		e.preventDefault();

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
						<p className="control has-icons-left">
							<input className="input" type="password" placeholder="Password"/>
							<span className="icon is-small is-left">
								<i className="fas fa-lock"/>
							</span>
							<p className="has-text-right mt1">
								<a className="has-text-right has-text-grey-dark" onClick={ () => props.history.push('reset-password') }>
									Forgot password?
								</a>
							</p>
						</p>
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
								  	<i className="fas fa-user"></i>
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

