<template>
<section class="row middle-xs mb3">
	<div class="col-xs-2 start-xs middle-xs ttu">
		<span class="f3">Users</span>
	</div>
	<div class="col-xs-10 start-xs middle-xs">
		<div v-if="!hasUsers">No users are defined. Add them in the main unRAID webGUI first.</div>
		<div v-else class="ttu">
			<UserButton v-for="user in users" :key="user.idx" :text="user.name" :selected="user.name === selectedUser" :on-click="selectUser.bind(this, user.name)">{{ user.name }}</UserButton>
		</div>
	</div>
</section>
</template>

<script lang="ts">
import Vue from 'vue'
import UserButton from './Button.vue'

import { SELECT_USER } from '../constants'

export default Vue.extend({
	name: 'users',

	components: { UserButton },

	computed: {
		users() {
			return this.$store.state.users
		},

		hasUsers() {
			return this.$store.state.users.length > 0
		},

		selectedUser() {
			return this.$store.state.selectedUser
		},
	},

	methods: {
		selectUser(name: string) {
			this.$store.commit(SELECT_USER, name)
		},
	},
})
</script>

<style scoped>

</style>
