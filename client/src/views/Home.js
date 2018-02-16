import React from 'react'

import { observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)
const banner = require('../img/iphone-hand.png')

export const Home = observer(({ store }) =>
	<div>
		<section className={cx('row')}>
			<div className={cx('col-xs-12', 'center-xs')}>
				<h2>ControlR plugin to generate Dockers/VMs user permissions</h2>
			</div>
		</section>

		<section className={cx('row')}>
			<div className={cx('col-xs-12')}>
				<img src={banner} alt="ControlR" />
			</div>
		</section>
	</div>,
)
