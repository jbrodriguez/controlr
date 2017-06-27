import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import { inject, observer } from 'mobx-react'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

@inject('store')
@observer
class App extends PureComponent {
	static propTypes = {
		store: PropTypes.object,
		app: PropTypes.object,
	}

	render() {
		const { app } = this.props
		return (
			<tr>
				<td className={cx('flex', 'flex-row', 'items-center')}>
					<div><img src={app.logo} alt="icon" className={cx('h2')} /></div>
					<div className={cx('ml2')} />
					<div> {app.name} </div>
				</td>
				<td className={cx('tc')}>
					<div className={cx('flex', 'items-center', 'justify-center')}>
						<input type="checkbox" checked={app.visible} onChange={this.onToggle(app.id, 'visible')} />
					</div>
				</td>
				<td className={cx('tc')}>
					<div className={cx('flex', 'items-center', 'justify-center')}>
						<input type="checkbox" checked={app.read} onChange={this.onToggle(app.id, 'read')} />
					</div>
				</td>
				<td className={cx('tc')}>
					<div className={cx('flex', 'items-center', 'justify-center')}>
						<input type="checkbox" checked={app.write} onChange={this.onToggle(app.id, 'write')} />
					</div>
				</td>
				<td className={cx('tc')}>
					<div className={cx('flex', 'items-center', 'justify-center')}>
						<input type="checkbox" checked={app.execute} onChange={this.onToggle(app.id, 'execute')} />
					</div>
				</td>
			</tr>
		)
	}

	onToggle = (id, perm) => e => {
		this.props.store.app.setPerm(id, perm)
	}
}

export const Apps = inject('store')(
	observer(({ store }) => {
		const apps = store.app.getApps()
		const hasApps = store.app.hasApps

		return (
			<section className={cx('row', 'middle-xs')}>
				<div className={cx('col-xs-12', 'ttu')}>
					<h2>Apps</h2>
				</div>

				<div className={cx('col-xs-12')}>
					<div>
						<strong>VISIBLE: </strong><span>If turned off, the docker/vm will not appear in the app</span>
						<br />
						<strong>READ: </strong><span>User can see the docker/vm, but actions are disabled</span><br />
						<strong>WRITE: </strong><span>User can edit/remove the docker/vm</span><br />
						<strong>EXECUTE: </strong><span>User can start/stop the docker/vm</span>
					</div>
				</div>

				<div className={cx('col-xs-12')}>
					{!hasApps
						? <div>No dockers/vms are available</div>
						: <table>
								<thead>
									<tr>
										<th className={cx('tl')} style={{ width: '500px' }}>APP</th>
										<th className={cx('tc')} style={{ width: '150px' }}>VISIBLE</th>
										<th className={cx('tc')} style={{ width: '150px' }}>READ</th>
										<th className={cx('tc')} style={{ width: '150px' }}>WRITE</th>
										<th className={cx('tc')} style={{ width: '150px' }}>EXEC</th>
									</tr>
								</thead>
								<tbody>
									{apps.map(app => <App key={app.id} app={app} />)}
								</tbody>
							</table>}
				</div>
			</section>
		)
	}),
)
