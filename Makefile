.PHONY: all start build test deps

all: deps test build

build-darwin-arm: test
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/darwin/ergo

build-linux-arm: test
	GOOS=linux GOARCH=arm64 go build -o bin/linux/ergo

build:
	@go build -o bin/ergo

start:
	@go run main.go run

test:
	@go test ./... -v

watch:
	find watch

deps:
	@go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v ergo | go get
