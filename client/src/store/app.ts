export default class App {
	public name: string
	public type: string
	public logo: string
	public id: string

	public visible: boolean
	public read: boolean
	public write: boolean
	public execute: boolean

	constructor(name: string, type: string, logo: string, id: string) {
		this.name = name
		this.type = type
		this.logo = logo
		this.id = id

		this.visible = false
		this.read = false
		this.write = false
		this.execute = false
	}

	public applyPerms = perms => {
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

	public setPerm = perm => {
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

	public clone = _ => {
		let app = new App(this.name, this.type, this.logo, this.id)

		app.applyPerms({
			visible: this.visible,
			allowed: { read: this.read, write: this.write, execute: this.execute },
		})

		return app
	}
}
