// Changing this will change the accepted minimal password length, and also the error messages when shown when inputting a too short password.
export const MIN_PASSWORD_LENGTH: number = 6;

export const checkForEmptyPassword = (password: string): boolean => {
	return !!password.length;
}

export const checkIfPasswordsMatch = (password: string, repeatPassword: string): boolean => {
	return password === repeatPassword;
}

export const checkForPasswordLength = (password: string): boolean => {
	return password.length >= MIN_PASSWORD_LENGTH;
}
