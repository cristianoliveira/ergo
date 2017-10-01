.PHONY: all start build test deps

all: deps bump-version build test

VERSION=`cat .version`
LDFLAGS_f1=-ldflags "-w -s -X main.VERSION=${VERSION}"

build-darwin-arm:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/darwin/ergo

build-linux-arm:
	GOOS=linux GOARCH=arm64 go build -o bin/linux/ergo

build-windows-i386:
	GOOS=windows GOARCH=386 go build -o bin/win/ergo.exe

bump-version:
	@git tag --sort=committerdate | tail -n 1 > .version
	cat .version

build: bump-version
	@go build ${LDFLAGS_f1} -o bin/ergo

start:
	@go run main.go run

test: build
	@go test ./... -v

watch:
	funzzy watch

deps:
	@go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v ergo | go get
