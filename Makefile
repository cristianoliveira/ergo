.PHONY: all start build test deps

all: deps test bump-version build

VERSION=`cat .version`
LDFLAGS_f1=-ldflags "-w -s -X main.VERSION=${VERSION}"

build-darwin-arm: test
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/darwin/ergo

build-linux-arm: test
	GOOS=linux GOARCH=arm64 go build -o bin/linux/ergo

bump-version:
	@git tag --sort=committerdate | tail -n 1 > .version
	cat .version

build: bump-version
	@go build ${LDFLAGS_f1} -o bin/ergo

start:
	@go run main.go run

test:
	@go test ./... -v

watch:
	funzzy watch

deps:
	@go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v ergo | go get
