import VueRouter from 'vue-router'
import Socket from './store/socket'

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

	users: IUsers
	userOrder: string[]

	apps: IUserApps
	appOrder: string[]

	qrcode: string
}

export interface IVersion {
	version: string
}

export interface IAuthParams {
	username: string
	password: string
	[key: string]: string
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

export interface IApp {
	name: string
	type: string
	logo: string
	id: string
	read: boolean
	write: boolean
	execute: boolean
	[perm: string]: boolean | string
}

export interface IApps {
	[appId: string]: IApp
}

export interface IUserApps {
	[userIdx: string]: IApps
}

export interface IUser {
	idx: string
	name: string
	desc: string
	perms: string
}

export interface IUsers {
	[idx: string]: IUser
}

export interface IPerm {
	read: boolean
	write: boolean
	exec: boolean
}

export interface IPerms {
	[name: string]: IPerm
}

export interface IRefreshed {
	users: IUser[]
	apps: IApp[]
}

export interface IEncode {
	[key: string]: string
}

export interface ISetPermArgs {
	user: string
	name: string
	perm: string
}
