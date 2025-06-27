<script setup>
import { ref, onMounted } from 'vue'
import { useStore } from '@/stores/store'
import { useWebSocketStore } from '@/stores/websocket'
const store = useStore()
const wsStore = useWebSocketStore()
const username = ref('')

onMounted(() => {
	const storedUsername = localStorage.getItem('userName')
	const storedUserId = localStorage.getItem('userID')
	if (storedUsername) username.value = storedUsername
})

async function login(e) {
	try {
		wsStore.connect()

		// Store in localStorage for persistence
		localStorage.setItem('userName', username.value)

	} catch (error) {
		console.error('Login failed:', error)
		alert('Failed to login. Please try again.')
	}
}

</script>

<template>
	<div v-if="!store.currentUser" class="login-container">
		<h1>Log in</h1>
		<form @submit.prevent="login">
			<input v-model="username" placeholder="Username" required />
			<button type="submit">Log In</button>
		</form>
	</div>
</template>
