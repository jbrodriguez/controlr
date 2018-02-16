<template>
<section class="row middle-xs mb3">
	<div class="col-xs-12 start-xs middle-xs ttu mb2">
		<span class="f3">Apps</span>
	</div>
	<div class="col-xs-12 start-xs middle-xs mb2">
		<div>
			<strong>VISIBLE: </strong><span>If turned off, the docker/vm will not appear in the app</span>
			<br />
			<strong>READ: </strong><span>User can see the docker/vm, but actions are disabled</span><br />
			<strong>WRITE: </strong><span>User can edit/remove the docker/vm</span><br />
			<strong>EXECUTE: </strong><span>User can start/stop the docker/vm</span>
		</div>
	</div>
	<div class="col-xs-12 start-xs middle-xs">
		<div v-if="!hasApps">No dockers/vms are available</div>
		<table v-else>
		<thead>
			<tr>
				<th class="tl w500">APP</th>
				<th class="tc w150">VISIBLE</th>
				<th class="tc w150">READ</th>
				<th class="tc w150">WRITE</th>
				<th class="tc w150">EXEC</th>
			</tr>
		</thead>
		<tbody>
			<tr v-for="app in apps" :key="app.id">
				<td class="flex flex-row items-center">
					<div><img :src="app.logo" alt="icon" class="h2" /></div>
					<div class="ml2" />
					<div>{{ app.name }}</div>
				</td>
				<td class="tc">
					<div class="flex items-center justify-center">
						<input type="checkbox" :checked="app.visible" @click="setPerm({id: app.id, perm: 'visible'})" />
					</div>
				</td>
				<td class="tc">
					<div class="flex' items-center justify-center">
						<input type="checkbox" :checked="app.read" @click="setPerm({id: app.id, perm: 'read'})"/>
					</div>
				</td>
				<td class="tc">
					<div class="flex items-center justify-center">
						<input type="checkbox" :checked="app.write" @click="setPerm({id: app.id, perm: 'write'})" />
					</div>
				</td>
				<td class="tc">
					<div class="flex items-center justify-center">
						<input type="checkbox" :checked="app.execute" @click="setPerm({id: app.id, perm: 'execute'})" />
					</div>
				</td>
			</tr>
		</tbody>
		</table>
	</div>
</section>
</template>

<script lang="ts">
import Vue from 'vue'
import Button from './Button.vue'

import { SET_PERM } from '../constants'

export default Vue.extend({
	name: 'apps',

	components: { Button },

	computed: {
		apps() {
			return this.$store.state.apps
		},

		hasApps() {
			return this.$store.state.apps.length > 0
		},

		selectedUser() {
			return this.$store.state.selectedUser
		},
	},

	methods: {
		setPerm(args) {
			this.$store.dispatch(SET_PERM, args)
		},
	},
})
</script>

<style scoped>
.w500 {
	width: 500px;
}

.w150 {
	width: 150px;
}
</style>
