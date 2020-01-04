#
# Makefile to perform "live code reloading" after changes to .go files.
#
# n.b. you must install fswatch (OS X: `brew install fswatch`)
#
# To start live reloading run the following command:
# $ make serve
#

mb_version := $(shell cat ../VERSION)
mb_count := $(shell git rev-list HEAD --count)
mb_hash := $(shell git rev-parse --short HEAD)
mb_plgver := $(shell cat ../PLGVER)

# binary name to kill/restart
PROG = controlr

# targets not associated with files
.PHONY: default build test coverage clean kill restart serve

# default targets to run when only running `make`
default: test

# clean up
clean:
	go clean

# run formatting tool and build
build: clean
	go build fmt
	go build -ldflags "-X main.Version=$(mb_plgver)|$(mb_version)" -gcflags "all=-trimpath=$GOPATH" -v -o controlr controlr.go

server: clean
	go build fmt
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(mb_plgver)|$(mb_version)" -gcflags "all=-trimpath=$GOPATH" -v -o ../dist/controlr controlr.go

# run unit tests with code coverage
test:
	go test -v

# generate code coverage report
coverage: test
	go build test -coverprofile=.coverage.out
	go build tool cover -html=.coverage.out

publish: client server
	rsync -avzP -e "ssh" dist/* $(SERVER):/usr/local/emhttp/plugins/controlr

deploy: server
	rsync -avzP -e "ssh" dist/* $(SERVER):/usr/local/emhttp/plugins/controlr
