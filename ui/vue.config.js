const ImageminPlugin = require('imagemin-webpack-plugin').default

module.exports = {
	devServer: {
		proxy: {
			'/': {
				target: 'http://lucy.apertoire.org:2378',
			},
		},
	},
	configureWebpack: config => {
		if (process.env.NODE_ENV === 'production') {
			// mutate config for production...
			return {
				plugins: [
					new ImageminPlugin({
						pngquant: {
							quality: '90-95',
						},
					}),
				],
			}
		}
	},
}
