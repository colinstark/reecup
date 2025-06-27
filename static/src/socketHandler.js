class SocketHandler {
	constructor() {
		this.connected = false
		this.socket = undefined
	}

	connect(socket_url, userId) {
		let url = new URL(socket_url || `ws://${window.location.host}/ws`)
		
		if (userId) {
			url.searchParams.append("userId", userId)
		}

		this.socket = new WebSocket(url)

		this.socket.onopen = () => {
			this.connected = true
			console.log('WebSocket connected')
			
			
		}

		this.socket.onclose = () => {
			this.connected = false
			console.log('WebSocket connection closed')
		}

		this.socket.onclose = (e) => handleMessage(e)

	}

	disconnect() {
		if (!this.socket) return

		this.socket.close()
		this.socket = null
		this.connected = false
	}

	get isConnected() {
		return this.connected
	}

	sendMessage(message) {
		if (!this.connected) {
			console.error('Cannot send message: WebSocket not connected')
			return
		}

		this.socket.send(JSON.stringify(message))
	}

}

function handleMessage(event) {
	const data = JSON.parse(event.data)

	if (data.type == "success") {
		console.log("Success")
		return
	}
	if (data.type == "error") {
		console.log("Error")
		return
	}

	console.log("got a message", data)
	switch (data.instruction) {
		case 'new_deck':
			gameStore.updateDeck(data.deck.stones)
			break
		case 'new_message':
			messageStore.addMessage(data.message)
			break
		case 'channel_created':
			channelStore.addChannel(data.channel)
			break
		case 'user_joined':
			// Update channel users if needed
			break
	}
}

export default SocketHandler
