import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useStore = defineStore('store', () => {
	const currentUser = ref(null)
	const users = ref({}) // map of userId -> user
	const games = ref({})

	function setCurrentUser(user) {
		currentUser.value = user
		addUser(user)
	}

	function addUser(user) {
		users.value[user.id] = user
	}

	function getUserById(id) {
		return users.value[id]
	}

	function setAllUsers(users) {
		return users.value = users
	}

	function setAllUsers(users) {
		return users.value = users
	}

	function setAllGames(newGames) {
		games.value = newGames
	}


	return {
		currentUser,
		users,
		setCurrentUser,
		games,
		setAllGames,
	}
})
