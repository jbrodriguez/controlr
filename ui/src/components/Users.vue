<template>
  <div>
    <section class="row middle-xs mb3">
      <div class="col-xs-2 start-xs middle-xs">
        <span class="f3">USERS</span>
      </div>
      <div class="col-xs-10 start-xs middle-xs">
        <div v-if="!hasUsers">No users are defined. Add them in the main unRAID webGUI first.</div>
        <div v-else>
          <UserButton
            v-for="user in users"
            :key="user.idx"
            :text="user.name"
            :selected="user.idx === selectedUser"
            :on-click="selectUser.bind(this, user.idx)"
          >{{ user.name }}</UserButton>
        </div>
      </div>
    </section>
    <section class="row middle-xs mb3">
      <div class="col-xs-12 start-xs middle-xs ttu">
        <a
          class="f6 link dim br2 ph3 pv2 mb2 dib white bg-dark-blue mr3"
          href="#"
          v-on:click="clearAll"
          v-if="selectedUser !== '' || isBusy"
        >CLEAR ALL</a>
        <span class="f6 br2 ba ph3 pv2 mb2 dib mid-gray mr3" v-else>CLEAR ALL</span>
        <a
          class="f6 link dim br2 ph3 pv2 mb2 dib white bg-dark-blue"
          href="#"
          v-on:click.prevent="saveChanges"
          v-if="isDirty && !isBusy"
        >SAVE</a>
        <span class="f6 br2 ba ph3 pv2 mb2 dib mid-gray" v-else>SAVE</span>
      </div>
    </section>
    <Apps :user="selectedUser"/>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import UserButton from './Button.vue'
import Apps from './Apps.vue'

import { IUser, IState } from '../types'
import { CLEAR_ALL, SAVE_CHANGES } from '../constants'

export default Vue.extend({
	name: 'users',

	data() {
		return {
			selectedUser: '',
		}
	},

	components: { UserButton, Apps },

	computed: {
		users(): IUser[] {
			const state: IState = this.$store.state
			return state.userOrder.map(idx => state.users[idx])
		},

		hasUsers(): boolean {
			return this.$store.state.userOrder.length > 0
		},

		isBusy(): boolean {
			return this.$store.state.isBusy
		},

		isDirty(): boolean {
			if (this.selectedUser === '') {
				return false
			}

			const state: IState = this.$store.state

			return (
				state.users[this.selectedUser].desc !==
				state.users[this.selectedUser].perms
			)
		},
	},

	methods: {
		selectUser(idx: string) {
			this.selectedUser = idx
		},

		clearAll() {
			this.$store.commit(CLEAR_ALL, this.selectedUser)
		},

		saveChanges() {
			this.$store.dispatch(SAVE_CHANGES, this.selectedUser)
		},
	},
})
</script>

<style scoped>
</style>
