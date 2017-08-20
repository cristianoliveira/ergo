.PHONY: all start build test deps

all: deps test build

build-darwin-arm: test
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/darwin/apitogo

build-linux-arm: test
	GOOS=linux GOARCH=arm64 go build -o bin/linux/apitogo

build:
	@go build -o bin/apitogo

start:
	@go run main.go run

test:
	@go test ./... -v

watch:
	find watch

deps:
	@go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v apitogo | go get
