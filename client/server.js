require('babel-core/register')({})

// var webpack = require('webpack');
// var WebpackDevServer = require('webpack-dev-server');
// var config = require('./webpack.config');
//
// new WebpackDevServer(webpack(config), {
// 	publicPath: config.output.publicPath,
// 	hot: true,
// 	historyApiFallback: true
// }).listen(3000, 'localhost', function (err, result) {
// 	if (err) {
// 	console.log(err);
// 	}
//
// 	console.log('Listening at localhost:3000');
// });


var path = require('path');
var express = require('express');
var webpack = require('webpack');
var webpackMiddleware = require('webpack-dev-middleware');
var webpackHotMiddleware = require('webpack-hot-middleware');
var config = require('./webpack.config.js');
var httpProxy = require('http-proxy');

// We need to add a configuration to our proxy server,
// as we are now proxying outside localhost
var proxy = httpProxy.createProxyServer({
  changeOrigin: true,
  ws: true,
})

const isDeveloping = process.env.NODE_ENV !== 'production'
// const port = isDeveloping ? 3000 : process.env.PORT;

const app = express()
const server = require('http').createServer(app)

if (isDeveloping) {
	const compiler = webpack(config);
	const middleware = webpackMiddleware(compiler, {
		publicPath: config.output.publicPath,
		contentBase: 'src',
		stats: {
			colors: true,
			hash: false,
			timings: true,
			chunks: false,
			chunkModules: false,
			modules: false
		}
	});

	app.use(middleware);
	app.use(webpackHotMiddleware(compiler));

	app.all('/api/*', function(req, res) {
		proxy.web(req, res, {target: 'http://blackbeard.apertoire.org:2378'})
	})

	app.all('/login', function(req, res) {
		proxy.web(req, res, {target: 'http://blackbeard.apertoire.org:2378'})
	})

	server.on('upgrade', function(req, socket, head) {
		proxy.ws(req, socket, head, {target: 'http://blackbeard.apertoire.org:2378'})
	});

	// app.all('/skt', function(req, res) {
	// 	proxy.ws(req, res, {
	// 		target: 'ws://wopr.apertoire.org:2378/skt'
	// 	})
	// })

	app.get('*', function response(req, res) {
		res.write(middleware.fileSystem.readFileSync(path.join(__dirname, 'dist/index.html')));
		res.end();
	});

} else {
	app.use(express.static(__dirname + '/dist'));
	app.get('*', function response(req, res) {
		res.sendFile(path.join(__dirname, 'dist/index.html'));
	});
}

const port = process.env.PORT || 3000

server.listen(port, '0.0.0.0', function onStart(err) {
	if (err) {
		console.log(err);
	}
	console.info('==> 🌎 Listening on port %s. Open up http://0.0.0.0:%s/ in your browser.', port, port);
});

// export default server

// // var server = require('http').createServer(app);
// server.listen(port, '0.0.0.0', function onStart(err) {
// 	if (err) {
// 		console.log(err);
// 	}
// 	console.info('==> 🌎 Listening on port %s. Open up http://0.0.0.0:%s/ in your browser.', port, port);
// });
