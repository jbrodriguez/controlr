import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import { observable, action } from 'mobx'
import { observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'
import views from '../config/views'

const cx = classNames.bind(styles)

@observer(['store'])
export default class Login extends PureComponent {
	@observable username = 'root'
	@observable password = ''

	static propTypes = {
		store: PropTypes.object,
	}

	render() {
		return (
			<form className={cx('flex', 'flex-row', 'items-center', 'ml2')} onSubmit={this.onSubmit}>
				<span>USER</span>
				<div className={cx('ml2')} />
				<input type="string" value={this.username} onChange={this.setUsername} />
				<div className={cx('ml2')} />
				<span>PASS</span>
				<div className={cx('ml2')} />
				<input type="password" value={this.password} onChange={this.setPassword} />
				<div className={cx('ml2')} />
				<input type="submit" value="Log In" />
			</form>
		)
	}

	@action('setUsername') setUsername = e => {
		this.username = e.target.value
	}

	@action('setPassword') setPassword = e => {
		this.password = e.target.value
	}

	onSubmit = e => {
		e.preventDefault()
		this.props.store.app.login({ username: this.username, password: this.password }, authenticated => {
			if (authenticated) {
				this.props.store.router.goTo(views.perms, {}, this.props.store)
			} else {
				this.password = ''
			}
		})
	}
}
