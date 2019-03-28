import { IAuthParams } from '../types'

const encode = (data: IAuthParams) => {
	const encoded = Object.keys(data).map(key => {
		const value = encodeURIComponent(data[key].toString())
		return `${key}=${value}`
	})
	return encoded.join('&')
}

const checkStatus = (response: any) => {
	if (response.ok) {
		return response
	} else {
		throw { code: response.status, msg: response.statusText }
	}
}

class Api {
	private host: string

	constructor() {
		this.host = document.location.origin
	}

	public getVersion: () => Promise<any> = () => {
		return fetch(this.host + '/version')
			.then(checkStatus)
			.then(resp => resp.json())
			.then(data => data)
	}

	public login: (params: IAuthParams) => Promise<any> = params => {
		return fetch(this.host + '/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
			body: encode(params),
		})
			.then(checkStatus)
			.then(resp => resp.json())
			.then(data => data)
	}
}

const api: Api = new Api()

export default api
