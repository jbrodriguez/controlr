import Vue from 'vue'
import Vuex, { Store } from 'vuex'
import { ActionTree, MutationTree, ActionContext } from 'vuex'

import * as constants from '../constants'
import {
	IState,
	ILogin,
	ILoginSuccess,
	IVersion,
	ILoginFailure,
	IUsers,
	IApps,
	IPerms,
	IPerm,
	IRefreshed,
	ISetPermArgs,
	IUserApps,
} from '@/types'

import api from './api'
import Socket from './socket'
import { parsePerms, createPerms } from '@/lib/utils'

Vue.use(Vuex)

const state: IState = {
	version: '',

	token: '',
	username: '',

	isAuthorized: false,
	isBusy: false,

	hasError: false,
	errCode: 0,
	errMsg: '',

	socket: null,

	users: {},
	userOrder: [],

	apps: {},
	appOrder: [],
}

const actions: ActionTree<IState, any> = {
	[constants.GET_VERSION]: async (context: ActionContext<IState, any>) => {
		const version = await api.getVersion()
		return context.commit(constants.GOT_VERSION, version)
	},

	[constants.LOGIN]: (context: ActionContext<IState, any>, { username, password, router }: ILogin) => {
		context.commit(constants.IS_BUSY, true)

		api.login({ username, password })
			.then(data => context.commit(constants.LOGIN_SUCCESS, { data, username, router }))
			.catch(err => context.commit(constants.LOGIN_FAILURE, err))
	},

	[constants.INIT]: (context: ActionContext<IState, any>) => {
		const proto = document.location.protocol === 'https:' ? 'wss' : 'ws'
		const host = `${proto}://${document.location.host}/skt/?token=${state.token}`

		const socket = new Socket(host, () => context.dispatch(constants.REFRESH))
		socket.receive((message: any) => {
			const packet = JSON.parse(message.data)
			if (typeof packet.topic === 'string' && packet.topic.length > 0) {
				context.dispatch(packet.topic, packet.payload)
			}
		})

		context.commit(constants.SET_SOCKET, socket)
	},

	[constants.REFRESH]: (context: ActionContext<IState, any>) => {
		context.commit(constants.IS_BUSY, true)

		if (!context.state.socket) {
			return
		}

		context.state.socket.send({ topic: constants.REFRESH, payload: {} })
	},

	[constants.REFRESHED]: (context: ActionContext<IState, any>, payload: IRefreshed) => {
		context.commit(constants.UPDATE_ALL, payload)
	},

	[constants.SAVE_CHANGES]: (context: ActionContext<IState, any>, idx: string) => {
		if (!context.state.socket) {
			return
		}

		const user = context.state.users[idx]

		context.commit(constants.IS_BUSY, true)
		context.state.socket.send({
			topic: constants.UPDATE_USER,
			payload: { idx: user.idx, name: user.name, perms: user.perms },
		})
	},

	[constants.USER_UPDATED]: async (context: ActionContext<IState, any>, { idx }) => {
		context.commit(constants.USER_UPDATED, idx)
		context.commit(constants.IS_BUSY, false)
	},

	[constants.ACCESS_ERROR]: (context: ActionContext<IState, any>, payload) => {
		context.commit(constants.SET_ERROR, payload)
	},
}

const mutations: MutationTree<IState> = {
	[constants.IS_BUSY]: (local: IState, isBusy: boolean) => (local.isBusy = isBusy),

	[constants.GOT_VERSION]: (local: IState, data: IVersion) => (local.version = data.version),

	[constants.LOGIN_SUCCESS]: (local: IState, { data, username, router }: ILoginSuccess) => {
		local.token = data.token
		local.isAuthorized = true
		local.isBusy = false
		local.username = username

		local.hasError = false
		local.errCode = 0
		local.errMsg = ''

		router.push(constants.PERMS)
	},

	[constants.LOGIN_FAILURE]: (local: IState, err: ILoginFailure) => {
		local.token = ''
		local.isAuthorized = false
		local.isBusy = false
		local.username = ''

		local.hasError = true
		local.errCode = err.code
		local.errMsg = err.msg
	},

	[constants.REMOVE_FEEDBACK]: (local: IState) => {
		local.hasError = false
		local.errCode = 0
		local.errMsg = ''
	},

	[constants.LOGOUT]: (local: IState) => {
		local.isBusy = false
		local.isAuthorized = false
		local.username = ''
		local.token = ''
		local.errCode = 0
		local.errMsg = ''
	},

	[constants.SET_SOCKET]: (local: IState, socket: Socket) => {
		local.socket = socket
	},

	[constants.UPDATE_ALL]: (local: IState, payload: IRefreshed) => {
		local.userOrder = payload.users.map(user => user.idx).sort()
		local.users = payload.users.reduce(
			(users, user) => {
				users[user.idx] = { ...user, perms: user.desc }
				return users
			},
			{} as IUsers,
		)

		local.appOrder = payload.apps.map(app => app.name.toLowerCase())
		local.apps = payload.users.reduce(
			(userApps, user) => {
				userApps[user.idx] = payload.apps.reduce(
					(apps, app) => {
						const appPerms = parsePerms(user.desc)
						const name = app.name.toLowerCase()

						// the app may not be present in the user desc.
						// in that case, we assign it default perms
						const perms = name in appPerms ? appPerms[name] : { read: false, write: false, exec: false }

						apps[name] = { ...app, ...perms }

						return apps
					},
					{} as IApps,
				)

				return userApps
			},
			{} as IUserApps,
		)

		// console.log(`users(${JSON.stringify(local.users)})-o`)
		// console.log(`apps(${JSON.stringify(local.apps)})-order(${JSON.stringify(local.appOrder)})`)

		local.isBusy = false
	},

	[constants.SET_PERM]: (local: IState, { user, name, perm }: ISetPermArgs) => {
		// apps = { "1": {"nzbget": { ..., read: true}} }
		const lcName = name.toLowerCase()
		local.apps[user][lcName][perm] = !local.apps[user][lcName][perm]

		const apps = local.appOrder.map(appName => local.apps[user][appName])

		local.users[user].perms = createPerms(apps)
	},

	[constants.CLEAR_ALL]: (local: IState, userIdx: string) => {
		local.users[userIdx].perms = ''

		for (const appName in local.apps[userIdx]) {
			if (!local.apps[userIdx].hasOwnProperty(appName)) {
				continue
			}

			local.apps[userIdx][appName].read = false
			local.apps[userIdx][appName].write = false
			local.apps[userIdx][appName].exec = false
		}
	},

	[constants.USER_UPDATED]: (local: IState, userIdx: string) => {
		local.users[userIdx].desc = local.users[userIdx].perms
	},

	[constants.SET_ERROR]: (local: IState, payload) => {
		local.isBusy = false

		local.hasError = true
		local.errCode = payload.errCode
		local.errMsg = payload.errMsg
	},
}

const store: Store<IState> = new Vuex.Store<IState>({
	state,
	mutations,
	actions,
})

export default store
