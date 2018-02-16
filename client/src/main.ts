import Vue from 'vue'
import { sync } from 'vuex-router-sync'

import 'flexboxgrid/css/flexboxgrid.min.css'
import 'tachyons/css/tachyons.min.css'

import App from './app.vue'
import router from './router'
import store from './store'
import * as constants from './constants'

sync(store, router)

const v = new Vue({
	el: '#app',
	router,
	store,
	render: h => h(App),
	created() {
		this.$store.dispatch(constants.GET_VERSION)
	},
})
