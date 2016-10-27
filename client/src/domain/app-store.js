import { observable, action, computed } from 'mobx'

import Api from '../lib/api'
import Socket from '../lib/socket'
import User from './user'
import App from './app'

export default class Store {
	api = new Api()
	socket = null

	@observable version = ''

	@observable username = ''
	@observable token = ''

	@observable isLoading = false
	@observable isBusy = false
	@observable isAuthorized = false

	@observable hasError = false
	@observable error = {}

	@observable users = []
	@observable apps = []

	@observable selectedUser = ''

	getUsers = _ => this.users.slice().sort( (a,b) => { return (a.idx > b.idx) ? 1 : ((b.idx > a.idx) ? -1 : 0)})
	getApps = _ => this.apps.slice()

	getVersion = _ => {
		this.api.getVersion()
			.then(action('gotVersion', data => {
				this.version = data.version
			}))
	}

	@action('login') login = (authParams, callback) => {
		this.isLoading = true

		this.api.login(authParams)
			.then(action('loginSuccess', data => {
				this.isLoading = false
				this.isAuthorized = true
				this.username = authParams.username
				this.token = data.token

				this.hasError = false
				this.error = {}

				callback(true)
			}))
			.catch(action('loginError', err => {
				this.isLoading = false
				this.isAuthorized = false
				this.username = ''
				this.token = ''

				this.hasError = true
				this.error = `${err.code} (${err.msg})`

				callback(false)
			}))
	}

	@action('logout') logout = _ => {
		this.isLoading = false
		this.isBusy = false
		this.isAuthorized = false
		this.username = ''
		this.token = ''
	}

	@action('removeFeedback') removeFeedback = _ => {
		this.hasError = false
		this.error = {}
	}

	@action('start') start = _ => {
		this.socket = new Socket(`ws://${document.location.host}/skt/?token=${this.token}`, this.load)
		this.socket.register('model/REFRESHED', this.refreshed)
		this.socket.register('model/USER_UPDATED', this.userUpdated)
		this.socket.register('model/ACCESS_ERROR', this.accessError)
	}

	@action('load') load = _ => {
		this.isBusy = true
		this.socket.send('model/REFRESH')
	}

	@action('model/REFRESHED') refreshed = payload => {
		this.users = payload.users.map( user => new User(user.idx, user.name, user.desc) )
		this.apps = payload.apps.map( app => new App(app.name, app.type, app.logo, app.id))

		if (payload.users.length > 0) {
			if (this.selectedUser !== '') {
				this.selectUser(this.selectedUser)
			} else {
				this.selectUser(payload.users[0].name)
			}
		}

		this.isBusy = false
	}

	@action('selectUser') selectUser = name => {
		const user = this.users.find( user => user.name === name)
		const perms = this.parsePerms(user.desc)

		this.apps.forEach( app => { app.applyPerms(perms[app.name.toLowerCase()]) })

		this.selectedUser = name
	}

	@action('setPerm') setPerm = (id, perm) => {
		const user = this.users.find( user => user.name === this.selectedUser)
		if (!user) {
			return
		}

		const index = this.apps.findIndex( app => app.id === id)

		let app = this.apps[index].clone()
		app.setPerm(perm)

		let newApps = [...this.apps]
		newApps[index] = app

		const perms = this.createPerms(newApps)

		this.isBusy = true
		this.socket.send('model/UPDATE_USER', {idx: user.idx, name: user.name, perms })
	}

	@action('model/USER_UPDATED') userUpdated = payload => {
		this.isBusy = false
		this.load()
	}

	@action('model/ACCESS_ERROR') accessError = payload => {
		this.isLoading = false
		this.isBusy = false

		this.hasError = true
		this.error = payload.error

	}

	@computed get hasUsers() { return this.users.length > 0 }
	@computed get hasApps() { return this.apps.length > 0 }

	createPerms = apps => {
		let perms = ''

		apps.forEach( app => {
			if (!app.visible) {
				return
			}

			if (perms === '') {
				perms = `${app.name.toLowerCase()}|${app.read ? 'r' : '-'}${app.write ? 'w' : '-'}${app.execute ? 'x' : '-'}`
			} else {
				perms += `,${app.name.toLowerCase()}|${app.read ? 'r' : '-'}${app.write ? 'w' : '-'}${app.execute ? 'x' : '-'}`
			}
		})

		return perms
	}

	parsePerms = raw => {
		// so raw is
		// nzbget:---|sonarr:rwx|plex:r-x|red 5:-wx
		// this gets me
		// ['nzbget:---', 'sonarr:rwx', 'plex:r-x', 'red 5:-wx']
		const items = raw.split(',')

		// console.log(`length-${items.length}`)
		let perms = {}

		for (let i = 0; i < items.length; i++) {
			// this gets me
			// ['nzbget', '---']
			const pair = items[i].split('|')

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
				perms[name] = {visible: false}
				continue
			}

			const r = perm.charAt(0)
			if (r === '-') {
				perms[name] = {visible: false}
				continue
			}

			const w = perm.charAt(1)
			const x = perm.charAt(2)

			if (w === 'w' && x === 'x') {
				perms[name] = {visible: true, allowed: {read: true, write: true, execute: true}}
				continue
			}

			let allowed = Object.assign({}, {read: true})

			if (w === 'w') {
				Object.assign(allowed, {write: true})
			}

			if (x === 'x') {
				Object.assign(allowed, {execute: true})
			}

			perms[name] = {visible: true, allowed}
		}

		return perms
	}
}
