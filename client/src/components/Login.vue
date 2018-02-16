<template>
<form class='flex flex-row items-center ml2' v-on:submit.prevent="submit">
	<span class="mr2">USER</span>
	<input class="mr2" type="string" v-model="username" />
	<span class="mr2">PASS</span>
	<input class="mr2" type="password" v-model="password" />
	<input type="submit" value="Log In" />
</form>
</template>

<script lang="ts">
import Vue from 'vue'
import * as constants from '../constants'

export default Vue.extend({
	name: 'login',

	data() {
		return {
			username: 'root',
			password: '',
		}
	},

	methods: {
		submit() {
			this.$store.dispatch(constants.LOGIN, {
				username: this.username,
				password: this.password,
				router: this.$router,
			})
			this.password = ''
		},
	},
})
</script>

<style scoped>

</style>
