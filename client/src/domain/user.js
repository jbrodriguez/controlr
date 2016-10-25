import { observable } from 'mobx'

export default class User {
	@observable idx
	@observable name
	@observable desc

	constructor(idx, name, desc) {
		this.idx = idx
		this.name = name
		this.desc = desc
	}
}
