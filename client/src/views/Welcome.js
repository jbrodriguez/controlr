import React from 'react'

import { observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const Welcome  = observer(['store'], ({ store, logout }) => (
	<div className={cx('ml2')}>
		<span className={cx('navy')}>welcome </span>
		<span className={cx('navy')}>{store.app.username} !</span>
		<span> (<a className={cx('dark-green')} href='/' onClick={logout}>log out</a>)</span>
	</div>
))

export default Welcome
