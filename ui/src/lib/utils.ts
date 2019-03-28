import { IPerms, IApp } from '@/types'

export const parsePerms = (raw: string) => {
	// so raw is
	// nzbget|---,sonarr|rwx,plex|r-x,redox 5|-wx
	// it's possible that an app is not present (invisible)
	// this gets me
	// ['nzbget|---', 'sonarr|rwx', 'plex|r-x', 'red 5|-wx']
	const apps = raw.split(',')

	return apps.reduce(
		(perms, app) => {
			// this gets me
			// ['nzbget', '---']
			const pair = app.split('|')

			// I can't parse it, so let's just ignore it
			// it's like it doesn't exist in the perm list
			if (pair.length !== 2) {
				return perms
			}

			const name = pair[0]
			const literals = pair[1]

			// user is sending some bogus perms, ignore this app
			if (literals.length !== 3) {
				return perms
			}

			// user is sending some bogus perms, ignore this app
			const match = /([r-][w-][x-])/.test(literals)
			if (!match) {
				return perms
			}

			perms[name] = {
				read: literals.charAt(0) === 'r',
				write: literals.charAt(1) === 'w',
				exec: literals.charAt(2) === 'x',
			}

			return perms
		},
		{} as IPerms,
	)
}

export const createPerms = (apps: IApp[]) =>
	apps.reduce((perms, app) => {
		if (!(app.read || app.write || app.exec)) {
			return perms
		}

		const prefix = perms === '' ? '' : ','

		return (
			perms +
			prefix +
			`${app.name.toLowerCase()}|${app.read ? 'r' : '-'}${app.write ? 'w' : '-'}${app.exec ? 'x' : '-'}`
		)
	}, '')

export default {
	parsePerms,
	createPerms,
}
