import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import { useStrict } from 'mobx'
import { observer } from 'mobx-react'
import { MobxRouter } from 'mobx-router'
import DevTools from 'mobx-react-devtools'
import classNames from 'classnames/bind'
import 'font-awesome-webpack'

import styles from './styles/core.scss'
// import views from './config/views'
import Feedback from './views/Feedback'
import Login from './views/Login'
import Welcome from './views/Welcome'

useStrict(true)
const cx = classNames.bind(styles)

const controlr = require('./img/controlr.png')
const vm = require('./img/v.png')
const unraid = require('./img/unraid.png')
const logo = require('./img/jbrio.png')
// const banner = require('./img/iphone-hand.png')
const apple = require('./img/appstore.png')
const google = require('./img/gbadge.png')

@observer(['store'])
export default class App extends PureComponent {
	static propTypes = {
		store: PropTypes.object,
	}

	componentDidMount() {
		this.props.store.app.getVersion()
	}

	render() {
		const { isAuthorized, isBusy, version } = this.props.store.app

		let alert = null
		if (this.props.store.app.hasError) {
			alert = <Feedback store={this.props.store.app} onRemoveFeedback={this.onRemoveFeedback} />
		}

		return (
			<div className={cx('container', 'flex', 'flex-column', 'body')}>

				<header className={cx('mb3', 'mt3')}>
					<nav className={cx('row')}>

						<ul className={cx('col-xs-12', 'col-sm-2')}>
							<li className={cx('center-xs', 'flex', 'bg-light-gray', 'pa1')}>
								<img src={controlr} />
							</li>
						</ul>

						<ul className={cx('col-xs-12', 'col-sm-10')}>

							<li className={cx(styles['bg-light-gray'], 'pa1')}>
								<section className={cx('row', 'middle-xs', 'items-center', 'justify-between')}>
									<div className={cx('col-xs-12', 'col-sm-10')}>
										{isAuthorized ? <Welcome logout={this.onLogout} /> : <Login />}
									</div>

									<div className={cx('col-xs-12', 'col-sm-1')}>
										{isBusy && <i className={cx('f3', 'dark-red', 'fa fa-refresh fa-spin')} />}
									</div>

									<div className={cx('col-xs-12', 'col-sm-1', 'middle-xs', 'end-xs', 'flex1')}>
										<a
											className={cx('flex1')}
											href="https://twitter.com/jbrodriguezio"
											title="@jbrodriguezio"
											target="_blank"
											rel="noopener noreferrer"
										>
											<i className={cx('f3', 'dark-green', 'fa fa-twitter')} />
										</a>
										<a
											className={cx('flex1', 'ml2')}
											href="https://github.com/jbrodriguez"
											title="github.com/jbrodriguez"
											target="_blank"
											rel="noopener noreferrer"
										>
											<i className={cx('f3', 'dark-green', 'fa fa-github')} />
										</a>
										<img src={vm} />
									</div>
								</section>

							</li>

						</ul>

					</nav>
				</header>

				<main className={cx('flex', 'flex-column')}>
					{alert}
					<MobxRouter />
				</main>

				<footer className={cx('bg-light-gray', 'pa2', 'mb3')}>
					<nav className={cx('row', 'legal', 'middle-xs')}>

						<ul className={cx('col-xs-12', 'col-sm-3')}>
							<div className={cx('flex', 'flex-row')}>
								<span className={cx('copyright', 'ml2')}>Copyright &copy; &nbsp;</span>
								<a className={cx('dark-green')} href="http://jbrodriguez.io/">Juan B. Rodriguez</a>
							</div>
						</ul>

						<ul className={cx('col-xs-12', 'col-sm-3', 'flex', 'end-xs')}>
							<span className={cx('version')}>{version !== '' ? 'v' + version : '-'}</span>
						</ul>

						<ul className={cx('col-xs-12', 'col-sm-6')}>
							<div className={cx('row', 'middle-xs', 'end-xs')}>
								<div className={cx('col-xs-12', 'middle-xs')}>
									<div className={cx('flex', 'end-xs')}>
										<a
											className={cx('flex', 'middle-xs')}
											href="https://itunes.apple.com/us/app/controlr/id1152586217?ls=1&mt=8"
										>
											<img alt="Download on the App Store" src={apple} />
										</a>
										<a
											className={cx('flex', 'middle-xs', 'ml2')}
											href="http://play.google.com/store/apps/details?id=io.jbrodriguez.controlr&utm_source=global_co&utm_medium=prtnr&utm_content=Mar2515&utm_campaign=PartBadge&pcampaignid=MKT-Other-global-all-co-prtnr-py-PartBadge-Mar2515-1"
										>
											<img alt="Get it on Google Play" src={google} />
										</a>
										<a
											className={cx('flex', 'middle-xs', 'ml2')}
											href="http://lime-technology.com/"
											title="Lime Technology"
											target="_blank"
											rel="noopener noreferrer"
										>
											<img src={unraid} alt="Logo for unRAID" />
										</a>
										<a
											className={cx('flex', 'middle-xs', 'ml2')}
											href="http://jbrodriguez.io/"
											title="jbrodriguez.io"
											target="_blank"
											rel="noopener noreferrer"
										>
											<img src={logo} alt="Logo for Juan B. Rodriguez" />
										</a>
									</div>
								</div>
							</div>
						</ul>

					</nav>

				</footer>
				<DevTools />
			</div>
		)
	}

	onRemoveFeedback = _ => {
		this.props.store.app.removeFeedback()
	}

	onLogout = _ => {
		this.props.store.app.logout()
	}
}
