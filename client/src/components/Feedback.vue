<template>
<section class="row mb3">
	<div class="col-xs-12">
		<div class="flex middle-xs between-sm bg-red ph2">
			<span class="navy">Unable to connect to server</span>
			<span v-on:click="remove" class="cursor">x</span>
		</div>
	</div>
	<div class="col-xs-12">
		<div class="bg-red ph2">
			<span class="white">({{ errCode }}) - {{ errMsg }}</span>
		</div>
	</div>
</section>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'
import { REMOVE_FEEDBACK } from '../constants'

export default Vue.extend({
	name: 'feedback',

	computed: mapState(['errCode', 'errMsg']),

	methods: {
		remove() {
			this.$store.commit(REMOVE_FEEDBACK)
		},
	},
})
</script>

<style scoped>
.cursor {
	cursor: pointer;
}
</style>
