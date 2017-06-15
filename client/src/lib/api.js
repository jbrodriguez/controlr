import fetch from 'isomorphic-fetch'

const encode = data => {
	const encoded = Object.keys(data).map(key => {
		const value = encodeURIComponent(data[key].toString())
		return `${key}=${value}`
	})
	return encoded.join('&')
}

const checkStatus = response => {
	if (response.ok) {
		return response
	} else {
		// eslint-disable-next-line
		throw { code: response.status, msg: response.statusText }
	}
}

export default class Api {
	constructor() {
		this.hostr = document.location.origin
	}

	getVersion = _ => {
		return fetch(this.hostr + '/version').then(checkStatus).then(resp => resp.json()).then(data => data)
	}

	login = authParams => {
		return fetch(this.hostr + '/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
			body: encode(authParams),
		})
			.then(checkStatus)
			.then(resp => resp.json())
			.then(data => data)
	}
}
