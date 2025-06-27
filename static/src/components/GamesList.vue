<script setup>
import { ref, onMounted } from 'vue'
import { useStore } from '@/stores/store'
const store = useStore()

</script>

<template>
	<section>
		<h1>Game list</h1>
		<div v-for="game in store.games">
			{{ game.id }}
		</div>
	</section>
</template>
