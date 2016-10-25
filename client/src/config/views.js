import React from 'react'

//models
import {Route} from 'mobx-router'

//components
import { Home } from '../views/Home'
import { Perms } from '../views/Perms'

const views = {
	home: new Route({ path: '/', component: <Home /> }),
	perms: new Route({
		path: '/permissions',
		component: <Perms />,
		beforeEnter: (route, params, store) => {
			if (!store.app.isAuthorized) {
				store.router.goTo(views.home, {}, store)
				return false
			}
		},
		onEnter: (route, params, store) => {
			store.app.start()
		}
	})
}

export default views
