export default class Socket {
	opened

	constructor(url, fn) {
		this.opened = fn

		this.ws = new WebSocket(url)
		this.ws.onopen = this.onOpen
		this.ws.onclose = this.onClose
		this.ws.onmessage = this.onMessage

		this.methods = {}
	}

	onOpen = () => {
		console.log('websocket open')
		this.opened && this.opened()
	}

	onClose = () => {
		console.log('websocket closed')
	}

	onMessage = reply => {
		const data = JSON.parse(reply.data)
		const fn = this.methods[data.topic]
		fn && fn(data.payload)
	}

	send = (topic, data) => {
		this.ws.send(JSON.stringify({ topic, payload: data }))
	}
	register = (action, fn) => {
		this.methods[action] = fn
	}
}
