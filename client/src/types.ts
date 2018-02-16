import VueRouter from 'vue-router'
import Socket from './store/socket'
import App from './store/app'

export interface IState {
	version: string

	username: string
	token: string

	isBusy: boolean
	isAuthorized: boolean

	hasError: boolean
	errCode: number
	errMsg: string

	socket: Socket | null

	users: IUser[]
	apps: App[]

	selectedUser: string
}

export interface IVersion {
	version: string
}

export interface IAuthParams {
	username: string
	password: string
}

export interface ILogin {
	username: string
	password: string
	router: VueRouter
}

export interface IToken {
	token: string
}

export interface ILoginSuccess {
	data: IToken
	username: string
	router: VueRouter
}

export interface ILoginFailure {
	code: number
	msg: string
}

export interface IUser {
	idx: string
	name: string
	desc: string
}
