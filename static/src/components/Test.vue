<script setup>
import { ref, onMounted, onBeforeMount } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'
import { useStore } from '@/stores/store'
import Name from '@/components/Name.vue'

const wsStore = useWebSocketStore()
const store = useStore()

const targetGame = ref('')

function newGame() {
	wsStore.sendMessage({
		instruction: "new_game"
	})
}

function getGames() {
	wsStore.sendMessage({
		instruction: "list_games"
	})
}

function deleteGame(id) {
	wsStore.sendMessage({
		instruction: "cancel_game",
		gameID: id
	})
}

function joinGame(id) {
	wsStore.sendMessage({
		instruction: "join_game",
		gameID: id
	})
}

function startGame(id) {
	wsStore.sendMessage({
		instruction: "start_game",
		gameID: id
	})
}

</script>

<template>
	<div>
		<h1>Testing ground</h1>
		<div>User ID: {{ store.currentUser.id }}</div>
		<div>User name: {{ store.currentUser.username }}</div>
		<Name />
		<button @click="newGame">New Game</button>
		<button @click="getGames">Fetch Games</button>

		<form>
			<input v-model="targetGame" placeholder="Target Game" required />
			<button type="button" @click.prevent="deleteGame(targetGame)">Delete Game</button>
			<button type="button" @click.prevent="joinGame(targetGame)">Join Game</button>
			<button type="button" @click.prevent="startGame(targetGame)">Start Game</button>
		</form>

	</div>
</template>
