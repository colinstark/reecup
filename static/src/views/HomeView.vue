<script setup>
import { computed } from 'vue'

import GamesList from '@/components/GamesList.vue'
import Deck from '@/components/Deck.vue'
import Test from '@/components/Test.vue'
import Login from '@/components/Login.vue'

import { useStore } from '@/stores/store'
const store = useStore()

const currentUser = computed({
	get() { return store.currentUser },
	set(val) { store.setCurrentUser(val) }
})

</script>

<template>
	<div v-if="!currentUser" class="container login">
		<Login />
	</div>
	<div v-else class="container">
		<Test />
		<GamesList />
		<Deck />
	</div>
</template>
