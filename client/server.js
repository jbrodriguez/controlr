const proxy = require('http-proxy-middleware')
const Bundler = require('parcel-bundler')
const express = require('express')
const process = require('process')

console.log(process.env.NODE_ENV)

let bundler = new Bundler('./index.html', { minify: false })
let app = express()

app.use('/', proxy({ target: 'http://wopr.apertoire.org:2378/', changeOrigin: true, secure: false }))
app.use(proxy('/skt', { target: 'ws://wopr.apertoire.org:2378/', changeOrigin: true, ws: true }))
app.use(bundler.middleware())

app.listen(Number(process.env.PORT || 1234))
