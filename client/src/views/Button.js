import React from 'react'

import { observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export const Button = observer(({text, selected, onClick}) => selected
	? <a className={cx('f6', 'link', 'br2', 'ph3', 'pv2', 'mb2', 'dib', 'white', 'bg-dark-blue', 'mr2')} href='#'>{text}</a>
	: <a className={cx('f6', 'link', 'dim', 'br2', 'ba', 'bw1', 'ph3', 'pv2', 'mb2', 'dib', 'dark-blue', 'mr2')} href='#' onClick={onClick}>{text}</a>
)
