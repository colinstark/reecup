<script setup>
import { useWebSocketStore } from '@/stores/websocket'
import { useStore } from '@/stores/store'
import { useGameStore } from '@/stores/game'

import Stone from "@/components/Stone.vue"

const store = useStore()
const gameStore = useGameStore()
const wsStore = useWebSocketStore()

function sendCursor(e) {
	wsStore.sendMessage({
		instruction: 'update_cursor',
		name: store.currentUser.username,
		id: store.currentUser.id,
		x: e.clientX,
		y: e.clientX
	})
}

function getDeck(e) {
	wsStore.sendMessage({
		instruction: 'get_deck'
	})
}

</script>

<template>
	<button @click="sendCursor">Send Cursor</button>
	<button @click="getDeck">Get a Deck</button>
	<div class="stones">
		<Stone v-for="stone in gameStore.game.deck" :color="stone.color" :face="stone.face" :joker="stone.joker" />
	</div>
</template>

<style scoped>
.stones {
	display: flex;
	flex-wrap: wrap;
	gap: 1rem;
}
</style>
