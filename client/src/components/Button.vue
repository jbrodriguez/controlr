<template>
<a v-if="selected" class="f6 link br2 ph3 pv2 dib white bg-dark-blue mr2" href="#">{{ text }}</a>
<a v-else class="f6 link dim br2 ba bw1 ph3 pv2 dib dark-blue mr2" href="#" v-on:click.prevent="onClick">{{ text }}</a>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
	name: 'user-button',
	props: ['text', 'selected', 'onClick'],
})
</script>

<style scoped>

</style>
