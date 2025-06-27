<script setup>
import { ref, onMounted } from 'vue'

import { useStore } from '@/stores/store'
import { useWebSocketStore } from '@/stores/websocket'
const store = useStore()
const wsStore = useWebSocketStore()

const username = ref('')

onMounted(() => {
	username.value = store.currentUser["username"]
})

async function login(e) {
	try {
		wsStore.sendMessage({
			"instruction": "update_name",
			"name": username.value,
		})
		localStorage.setItem('userName', username.value)

	} catch (error) {
		console.error('Login failed:', error)
		alert('Failed to login. Please try again.')
	}
}


</script>


<template>
	<form @submit.prevent="login">
		<input v-model="username" placeholder="Username" required />
		<button type="submit">Update</button>
	</form>
</template>
