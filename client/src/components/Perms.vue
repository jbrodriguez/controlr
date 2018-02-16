<template>
<section class="row">
	<div class="col-xs-12 middle-xs center-xs">
		<Users />
		<Apps />
	</div>
</section>
</template>

<script lang="ts">
import Vue from 'vue'

import Users from './Users.vue'
import Apps from './Apps.vue'

import { INIT } from '../constants'

export default Vue.extend({
	name: 'perms',

	components: { Users, Apps },

	created() {
		this.$store.dispatch(INIT)
	},
})
</script>

<style scoped>

</style>
