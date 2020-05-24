export interface IUserCreate {
	email: string;
	password: string;
}

export interface IUser {
	id:         string;
	email:      string;
	isAdmin:    boolean;
	created_at: Date;
	updated_at: Date;
}
