const proxy = require('http-proxy-middleware')
const Bundler = require('parcel-bundler')
const express = require('express')

// let bundler = new Bundler('index.html')
let app = express()

app.use('/', proxy({ target: 'https://lucy.apertoire.org:2378/', changeOrigin: true, secure: false }))
// app.use('/skt', proxy({ target: 'https://lucy.apertoire.org:2378/', changeOrigin: true, ws: true }))

// app.use(bundler.middleware())

app.listen(Number(process.env.PORT || 1234))
