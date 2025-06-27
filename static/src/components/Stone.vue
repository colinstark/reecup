<script setup>
	defineProps({
		color: String,
		face: String,
		joker: Boolean
	})

	function randomRotate() {
		let amount = -2 + Math.random() * 4
		return `${amount}deg`
	}
</script>

<template>
	<div class="stone" :class="color" :style="{rotate: randomRotate()}">
		<div class="indent" v-if="face || joker">
			<div v-if="joker" class="joker">♕</div>
			<div v-if="!joker">{{face}}</div>
		</div>
	</div>
</template>

<style scoped>
	.stone {
		width: 4rem;
		height: 5.5rem;
		border-radius: 0.5rem;
		display:flex;
		align-items: center;
		justify-content: center;
		padding: 0.5rem 0.5rem 2rem;
		background: rgb(255,255,240);
		color: black;
		font-size: 1.5rem;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2), inset 0 -8px 0 rgba(0, 0, 0, 0.1);
		filter: brightness(0.98);
		transition: all .1s ease-out;
		&.red {
			color: red;
		}
		&.blue {
			color: blue;
		}
		&.yellow {
			color: #FF9900;
		}

		&:hover {
			scale: 1.05;
			translate: 0 -5%;
			filter: brightness(1.02);
		}
	}
	.indent {
		background: linear-gradient(0deg, rgba(0,0,0,0.05), rgba(0,0,0,0));
		width: 3rem;
		height: 3rem;
		display:flex;
		align-items: center;
		justify-content: center;
		border-radius: 100%;
		& > div {
			cursor: default;
			user-select: none !important;
			-webkit-user-select: none !important;
			font-weight: bold !important;
		}

		& > .joker {
			scale: 2;
		}
	}
</style>
