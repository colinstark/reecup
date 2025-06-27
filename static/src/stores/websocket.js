// src/stores/websocket.js
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useStore } from './store'

export const useWebSocketStore = defineStore('websocket', () => {
	const socket = ref(null)
	const connected = ref(false)

	const store = useStore()

	function connect(userId) {
		if (socket.value) {
			socket.value.close()
		}
		let url = new URL(`ws://${window.location.host}`)
		url.pathname = "/ws"

		// if in vite dev mode, manually set the correct ws port
		if (import.meta.env.DEV) {
			url.port = "8080"
		}

		if (userId) {
			url.searchParams.append("userID", userId)
		}

		socket.value = new WebSocket(url)

		socket.value.onopen = () => {
			connected.value = true
			console.log('WebSocket connected')
			if (userId) {
				let userName = localStorage.getItem("userName")
				console.log('Sending name..', userName)
				sendMessage({
					"instruction": "update_name",
					"name": userName,
				})
			}
		}

		socket.value.onclose = () => {
			connected.value = false
			console.log('WebSocket disconnected')
		}

		socket.value.onmessage = (event) => {
			const data = JSON.parse(event.data)

			if (data.type == "user_id_commissioned") {
				localStorage.setItem("userID", data.user_id)
				let userName = localStorage.getItem("userName")

				store.setCurrentUser({
					id: data.user_id,
					username: userName
				})

				console.log("Reconnecting...")
				connect(data.user_id)
				return
			}
			if (data.type == "error") {
				console.log("Error")
				return
			}

			console.log("got a message", data)
			switch (data.instruction) {
				case 'name_updated':
					store.currentUser.username = data.name
					break
				case 'new_deck':
					store.updateDeck(data.deck.stones)
					break
				case 'games_list':
					console.log("GAMES LIST", data.games)
					store.setAllGames(data.games)
					break
				case 'game_created':
					console.log("GAME GOT CREATED")
					// gameStore.updateDeck(data.deck.stones)
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
	}
	function sendMessage(message) {
		if (!connected.value) {
			console.error('Cannot send message: WebSocket not connected')
			return
		}

		socket.value.send(JSON.stringify(message))
	}

	function disconnect() {
		if (socket.value) {
			socket.value.close()
			socket.value = null
			connected.value = false
		}
	}

	return {
		connected,
		connect,
		sendMessage,
		disconnect
	}
})
