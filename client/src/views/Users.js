import React, { Component } from 'react'

import { observable, action } from 'mobx'
import { observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'
import { Button } from './Button'

const cx = classNames.bind(styles)

@observer(['store'])
export class Users extends Component {
	render() {
		const users = this.props.store.app.getUsers()
		const hasUsers = this.props.store.app.hasUsers
		const selected = this.props.store.app.selectedUser

		return (
		<section className={cx('row', 'middle-xs')}>
			<div className={cx('col-xs-1', 'ttu')}>
				<h2>Users</h2>
			</div>

			<div className={cx('col-xs-11')}>
			{ !hasUsers
				? <div>No users are defined. Add them in the main unRAID webGUI first.</div>
				: <div className={cx('ttu')}>
					{users.map( user => <Button key={user.name} text={user.name} selected={user.name === selected} onClick={this.onClick(user.name)}/> )}
				</div>
			}
			</div>
		</section>
		)
	}

	onClick = name => e => {
		e.preventDefault()
		this.props.store.app.selectUser(name)
	}
}
