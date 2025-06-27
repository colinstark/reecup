<script setup>
import { onMounted } from 'vue'
import { useStore } from '@/stores/store'
import { useWebSocketStore } from '@/stores/websocket'
import { RouterLink, RouterView } from 'vue-router'

const store = useStore()
const wsStore = useWebSocketStore()

onMounted(() => {
	// Check for existing user session
	const storedUserId = localStorage.getItem('userID')
	const storedUsername = localStorage.getItem('userName')

	if (storedUserId && storedUsername) {
		store.setCurrentUser({
			id: storedUserId,
			username: storedUsername
		})

		// Initialize websocket connection
		wsStore.connect(store.currentUser.id)

		// Load user channels
		// channelStore.fetchUserChannels(userId)
	}
})
</script>

<template>
	<header>
		<div class="wrapper">
			<nav>
				<RouterLink to="/">Home</RouterLink>
				<RouterLink to="/game">Game</RouterLink>
			</nav>
		</div>
	</header>

	<RouterView />

</template>

<style scoped>
header {
	line-height: 1.5;
	max-height: 100vh;
}

main {
	max-height: 100vh;
}

.logo {
	display: block;
	margin: 0 auto 2rem;
}

nav {
	width: 100%;
	font-size: 12px;
	text-align: center;
	margin-top: 2rem;
	display: flex;
}

nav a.router-link-exact-active {
	color: var(--color-text);
}

nav a.router-link-exact-active:hover {
	background-color: transparent;
}

nav a {
	display: inline-block;
	padding: 0 1rem;
	border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
	border: 0;
}

@media (min-width: 1024px) {}
</style>
