import Vue from 'vue'
import Vuex, { Store } from 'vuex'
import { ActionTree, MutationTree, ActionContext, Action } from 'vuex'

import * as constants from '../constants'
import { IState, ILogin, ILoginSuccess, IVersion, ILoginFailure } from '../types'

import api from './api'
import Socket from './socket'
import App from './app'

Vue.use(Vuex)

const parsePerms = raw => {
	// so raw is
	// nzbget:---|sonarr:rwx|plex:r-x|red 5:-wx
	// this gets me
	// ['nzbget:---', 'sonarr:rwx', 'plex:r-x', 'red 5:-wx']
	const items = raw.split(',')

	const perms = {}

	for (const item of items) {
		// this gets me
		// ['nzbget', '---']
		const pair = item.split('|')

		// I can't parse it, so let's just forget about it
		// it's like it doesn't exist in the perm list
		if (pair.length !== 2) {
			continue
		}

		const name = pair[0].toLowerCase()
		const perm = pair[1]

		// you're sending some bogus perms, forget about it
		if (perm.length !== 3) {
			continue
		}

		const match = /([r-][w-][x-])/.test(perm)
		if (!match) {
			perms[name] = { visible: false }
			continue
		}

		const r = perm.charAt(0)
		if (r === '-') {
			perms[name] = { visible: false }
			continue
		}

		const w = perm.charAt(1)
		const x = perm.charAt(2)

		if (w === 'w' && x === 'x') {
			perms[name] = { visible: true, allowed: { read: true, write: true, execute: true } }
			continue
		}

		const allowed = {
			read: true,
			write: w === 'w',
			execute: x === 'x',
		}

		perms[name] = { visible: true, allowed }
	}

	return perms
}

const createPerms = apps => {
	let perms = ''

	apps.forEach(app => {
		if (!app.visible) {
			return
		}

		if (perms === '') {
			perms = `${app.name.toLowerCase()}|${app.read ? 'r' : '-'}${app.write ? 'w' : '-'}${
				app.execute ? 'x' : '-'
			}`
		} else {
			perms += `,${app.name.toLowerCase()}|${app.read ? 'r' : '-'}${app.write ? 'w' : '-'}${
				app.execute ? 'x' : '-'
			}`
		}
	})

	return perms
}

const sleep = ms => {
	return new Promise(resolve => setTimeout(resolve, ms))
}

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

	users: [],
	apps: [],

	selectedUser: '',
}

const actions: ActionTree<IState, any> = {
	[constants.GET_VERSION]: (store: ActionContext<IState, any>) => {
		return api.getVersion().then(version => store.commit(constants.GOT_VERSION, version))
	},

	[constants.LOGIN]: (store: ActionContext<IState, any>, { username, password, router }: ILogin) => {
		store.commit(constants.IS_BUSY, true)

		api
			.login({ username, password })
			.then(data => store.commit(constants.LOGIN_SUCCESS, { data, username, router }))
			.catch(err => store.commit(constants.LOGIN_FAILURE, err))
	},

	[constants.INIT]: (store: ActionContext<IState, any>) => {
		const proto = document.location.protocol === 'https:' ? 'wss' : 'ws'
		const host = `${proto}://${document.location.host}/skt/?token=${state.token}`

		const socket = new Socket(host, () => store.dispatch(constants.REFRESH))
		socket.receive(message => {
			const packet = JSON.parse(message.data)
			if (typeof packet.topic === 'string' && packet.topic.length > 0) {
				store.dispatch(packet.topic, packet.payload)
			}
		})

		store.commit(constants.SET_SOCKET, socket)
	},

	[constants.REFRESH]: (store: ActionContext<IState, any>) => {
		store.commit(constants.IS_BUSY, true)
		store.state.socket.send({ topic: constants.REFRESH, payload: {} })
	},

	[constants.REFRESHED]: (store: ActionContext<IState, any>, payload) => {
		store.commit(constants.UPDATE_ALL, payload)
	},

	[constants.SET_PERM]: (store: ActionContext<IState, any>, { id, perm }) => {
		const user = store.state.users.find(user => user.name === store.state.selectedUser)
		if (!user) {
			return
		}

		const app = store.state.apps.find(app => app.id === id)

		app.setPerm(perm)

		const perms = createPerms(store.state.apps)

		store.commit(constants.IS_BUSY, true)
		store.state.socket.send({ topic: constants.UPDATE_PERMS, payload: { idx: user.idx, name: user.name, perms } })
	},

	[constants.PERMS_UPDATED]: async (store: ActionContext<IState, any>) => {
		await sleep(1000)
		store.state.socket.send({ topic: constants.REFRESH, payload: {} })
	},

	[constants.ACCESS_ERROR]: (store: ActionContext<IState, any>, payload) => {
		store.commit(constants.SET_ERROR, payload)
	},
}

const mutations: MutationTree<IState> = {
	[constants.IS_BUSY]: (state: IState, isBusy: boolean) => (state.isBusy = isBusy),

	[constants.GOT_VERSION]: (state: IState, data: IVersion) => (state.version = data.version),

	[constants.LOGIN_SUCCESS]: (state: IState, { data, username, router }: ILoginSuccess) => {
		state.token = data.token
		state.isAuthorized = true
		state.isBusy = false
		state.username = username

		state.hasError = false
		state.errCode = 0
		state.errMsg = ''

		router.push(constants.PERMS)
	},

	[constants.LOGIN_FAILURE]: (state: IState, err: ILoginFailure) => {
		state.token = ''
		state.isAuthorized = false
		state.isBusy = false
		state.username = ''

		state.hasError = true
		state.errCode = err.code
		state.errMsg = err.msg
	},

	[constants.REMOVE_FEEDBACK]: (state: IState) => {
		state.hasError = false
		state.errCode = 0
		state.errMsg = ''
	},

	[constants.LOGOUT]: (state: IState) => {
		state.isBusy = false
		state.isAuthorized = false
		state.username = ''
		state.token = ''
		state.errCode = 0
		state.errMsg = ''
	},

	[constants.SET_SOCKET]: (state: IState, socket: Socket) => {
		state.socket = socket
	},

	[constants.UPDATE_ALL]: (state: IState, payload) => {
		state.users = payload.users
		state.apps = payload.apps.map(app => new App(app.name, app.type, app.logo, app.id))

		if (payload.users.length > 0) {
			let name = ''

			if (state.selectedUser !== '') {
				name = state.selectedUser
			} else {
				name = payload.users[0].name
			}

			const user = state.users.find(user => user.name === name)
			const perms = parsePerms(user.desc)

			state.apps.forEach(app => {
				app.applyPerms(perms[app.name.toLowerCase()])
			})

			state.selectedUser = name
		}

		state.isBusy = false
	},

	[constants.SELECT_USER]: (state: IState, name: string) => {
		const user = state.users.find(user => user.name === name)
		const perms = parsePerms(user.desc)

		state.apps.forEach(app => {
			app.applyPerms(perms[app.name.toLowerCase()])
		})

		state.selectedUser = name
	},

	[constants.SET_ERROR]: (state: IState, payload) => {
		state.isBusy = false

		state.hasError = true
		state.errCode = payload.errCode
		state.errMsg = payload.errMsg
	},
}

const store: Store<IState> = new Vuex.Store<IState>({
	state,
	mutations,
	actions,
})

export default store
