<template>
<nav class="row pa1 middle-xs">
	<ul class="col-xs-12 col-sm-4">
		<li class="row start-xs">
			<div class="col-xs-12">
				<span> Copyright &copy;</span>
				<a class="dark-green" href="https://jbrio.net/posts/">Juan B. Rodriguez</a>
			</div>
		</li>
	</ul>

	<ul class="col-xs-12 col-sm-3">
		<li class="middle-xs end-xs">
			<span>{{ version }}</span>
		</li>
	</ul>

	<ul class="col-xs-12 col-sm-5">
		<li class="row middle-xs end-xs">
			<div class="col-xs-12">
			<div class="flex end-xs">
				<a class="flex end-xs" title="Download on the App Store" href="https://itunes.apple.com/us/app/controlr/id1152586217?ls=1&mt=8" target="_blank" rel="noopener noreferrer">
					<img src="../assets/appstore.png"/>
				</a>

				<a class="flex middle-xs ml2" title="Get it on Google Play" href="https://play.google.com/store/apps/details?id=io.jbrodriguez.controlr&utm_source=global_co&utm_medium=prtnr&utm_content=Mar2515&utm_campaign=PartBadge&pcampaignid=MKT-Other-global-all-co-prtnr-py-PartBadge-Mar2515-1" target="_blank" rel="noopener noreferrer">
					<img src="../assets/gbadge.png"/>
				</a>

				<a class="flex middle-xs ml2" title="Lime Technology" href="https://lime-technology.com" target="_blank" rel="noopener noreferrer">
					<img src="../assets/unraid.png"/>
				</a>
				
				<a class="flex middle-xs ml2" title="Juan B. Rodriguez" href="https://jbrio.net/posts/" target="_blank" rel="noopener noreferrer">
					<img src="../assets/jbrio.png"/>
				</a>
			</div>
			</div>
		</li>
	</ul>
</nav>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

export default Vue.extend({
	name: 'footer-bar',
	computed: mapState(['version']),
})
</script>

<style scoped>

</style>
