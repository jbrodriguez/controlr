import React from 'react'

import { observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const Feedback = observer(({ store, onRemoveFeedback }) =>
	<section className={cx('row', 'mb3')}>
		<div className={cx('col-xs-12')}>
			<div className={cx('flex', 'middle-xs', 'between-sm', 'bg-red', 'ph2')}>
				<span className={cx('navy')}>Unable to connect to server</span>
				<i className={cx('fa fa-remove', 'navy')} onClick={onRemoveFeedback} />
			</div>
		</div>
		<div className={cx('col-xs-12')}>
			<div className={cx('bg-red', 'ph2')}>
				<span className={cx('white')}>{store.error} </span>
			</div>
		</div>
	</section>,
)

export default Feedback
