export default class Socket {
	private skt: WebSocket

	constructor(url: string, fn: () => void) {
		this.skt = new WebSocket(url)

		this.skt.onopen = () => {
			// console.log('Connection opened')
			if (!fn) {
				return
			}

			fn()
		}
		// this.skt.onclose = () => console.log('Connection is closed...')
	}

	public receive(fn: any) {
		this.skt.onmessage = fn
	}

	public send(packet: any) {
		// packet = { topic, payload }
		this.skt.send(JSON.stringify(packet))
	}
}
