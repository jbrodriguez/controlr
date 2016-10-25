var path = require('path');
var HtmlWebpackPlugin = require('html-webpack-plugin')
var webpack = require('webpack');

// new HtmlWebpackPlugin({
// 	template: 'index.tpl.html',
// 	inject: 'body',
// 	filename: 'index.html'
// }),

		// 'react-hot-loader/patch',
		// 'webpack-dev-server/client?http://localhost:3000',
		// 'webpack/hot/only-dev-server',
		// './src/index'

		// output: {
		// 	path: path.join(__dirname, 'dist'),
		// 	filename: 'bundle.js',
		// 	publicPath: '/static/'
		// },


module.exports = {
	devtool: 'eval',
	entry: [
		'webpack-hot-middleware/client?reload=true',
		'./src/index.js'
	],
	output: {
		path: path.join(__dirname, '/dist/'),
		filename: '[name].js',
		publicPath: '/'
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: 'index.tpl.html',
			inject: 'body',
			filename: 'index.html'
		}),
		new webpack.optimize.OccurenceOrderPlugin(),
		new webpack.HotModuleReplacementPlugin(),
	],
	resolve: {
		extensions: ['', '.js', '.jsx']
	},
	module: {
		loaders: [{
			test: /\.jsx?$/,
			loaders: ['babel'],
			include: path.join(__dirname, 'src')
		}, {
			test: /\.scss$/,
			include: path.join(__dirname, 'src/styles'),
			loaders: [
				'style',
				'css?modules&localIdentName=[name]---[local]---[hash:base64:5]',
				'postcss',
				'sass'
			]
		}, {
			test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
			loader: "url-loader?limit=10000&minetype=application/font-woff&name=img/[name]-[hash:7].[ext]"
		}, {
			test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
			loader: "file?hash=sha512&digest=hex&name=img/[name]-[hash:7].[ext]"
		}, {
    		test: /\.(jpe?g|png|gif|svg)$/i,
			include: path.resolve(__dirname, 'src/img'),
			loaders: [
				'file?hash=sha512&digest=hex&name=img/[name]-[hash:7].[ext]',
				'image-webpack?{progressive:true, optimizationLevel: 7, interlaced: false, pngquant:{quality: "65-90", speed: 4}}'
			]
		}, {
			test: /\.css$/,
			loader: 'style!css?modules&localIdentName=[name]---[local]---[hash:base64:5]'
		}]
	}
};
