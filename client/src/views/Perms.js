import React from 'react'

import { inject, observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'
import { Users } from './Users'
import { Apps } from './Apps'

const cx = classNames.bind(styles)

export const Perms = inject('store')(
	observer(({ store }) => {
		return store.app.isLoading
			? <div>Loading ...</div>
			: <div>
					<Users />
					<div className={cx('mb3')} />
					<Apps />
				</div>
	}),
)
