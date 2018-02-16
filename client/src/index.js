import React from 'react'
import { render } from 'react-dom'

import { startRouter } from 'mobx-router'
import { Provider } from 'mobx-react'

import store from './domain/store'
import views from './config/views'
import App from './app'

startRouter(views, store)

render(
	<Provider store={store}>
		<App />
	</Provider>,
	document.getElementById('mnt')
)
