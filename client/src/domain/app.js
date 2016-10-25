import { observable, action } from 'mobx'

export default class App {
	@observable name
	@observable type
	@observable logo
	@observable id

	@observable visible
	@observable read
	@observable write
	@observable execute

	constructor(name, type, logo, id) {
		this.name = name
		this.type = type
		this.logo = logo
		this.id = id
	}

	@action('applyPerms') applyPerms = perms => {
		if (perms) {
			this.visible = perms.visible || false
			this.read = perms.allowed ? perms.allowed.read || false : false
			this.write = perms.allowed ? perms.allowed.write || false : false
			this.execute = perms.allowed ? perms.allowed.execute || false : false
		} else {
			this.visible = false
			this.read = false
			this.write = false
			this.execute = false
		}
	}

	@action('setPerm') setPerm = perm => {
		if (perm === 'visible') {
			if (this.visible) {
				this.read = false
				this.write = false
				this.execute = false
			} else {
				this.read = true
			}

			this.visible = !this.visible

			return
		}

		if (perm === 'read' && this.read) {
			this.visible = false
			this.read = false
			this.write = false
			this.execute = false

			return
		}

		this.visible = true
		this[perm] = !this[perm]
	}

	clone = _ => {
		let app = new App(this.name, this.type, this.logo, this.id)
		app.applyPerms({visible: this.visible, allowed: {read: this.read, write: this.write, execute: this.execute}})
		return app
	}
}
