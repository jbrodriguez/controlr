import Vue from 'vue'
import VueRouter from 'vue-router'

import { HOME, PERMS } from '../constants'

import Home from '../components/Home.vue'
import Perms from '../components/Perms.vue'

import store from '../store'

Vue.use(VueRouter)

export default new VueRouter({
	routes: [
		{ path: HOME, component: Home },
		{
			path: PERMS,
			component: Perms,
			beforeEnter: (to, from, next) => {
				if (!store.state.isAuthorized) {
					next(HOME)
				} else {
					next()
				}
			},
		},
	],
})
