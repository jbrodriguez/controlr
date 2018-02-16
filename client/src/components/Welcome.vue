<template>
<section class="ml2">
	<span class="navy">welcome </span>
	<span class='navy'>{{ username }} ! </span>
	<span>(<a class="dark-green" href="/" @click="logout">log out</a>)</span>
</section>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'
import { LOGOUT } from '../constants'

export default Vue.extend({
	name: 'welcome',

	methods: {
		logout() {
			this.$store.commit(LOGOUT)
		},
	},

	computed: mapState(['username']),
})
</script>

<style scoped>

</style>
