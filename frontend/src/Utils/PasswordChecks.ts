export const checkForEmptyPassword = (password: string): boolean => {
	return !!password.length;
}

export const checkIfPasswordsMatch = (password: string, repeatPassword: string): boolean => {
	return password === repeatPassword;
}

export const checkForPasswordLength = (password: string): boolean => {
	return password.length >= 8;
}
