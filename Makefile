NO_COLOR=\33[0m

OK_COLOR=\33[32m

ERROR_COLOR=\033[31m

WARN_COLOR=\33[33m

.PHONY: all start build test test-integration deps help fmt vet lint tools

VERSION=`cat .version`
LDFLAGS_f1=-ldflags "-w -s -X main.VERSION=${VERSION}"

help:
	@(echo "${WARN_COLOR}Usage:${NO_COLOR}")
	@(echo "${OK_COLOR}make all${NO_COLOR}                      Run the tests and build the executable")
	@(echo "${OK_COLOR}make help${NO_COLOR}                     Show this help")
	@(echo "${OK_COLOR}make build-darwin-arm${NO_COLOR}         Builds the executable for osx arm")
	@(echo "${OK_COLOR}make build-linux-arm${NO_COLOR}          Builds the executable for linux arm")
	@(echo "${OK_COLOR}make build-windows-i386${NO_COLOR}       Builds the executable for windows")
	@(echo "${OK_COLOR}make bump-version${NO_COLOR}             Write the new version for latter use")
	@(echo "${OK_COLOR}make build${NO_COLOR}                    Builds the executable for current system")
	@(echo "${OK_COLOR}make start${NO_COLOR}                    Starts ergo")
	@(echo "${OK_COLOR}make fmt${NO_COLOR}                      Run gofmt on the source code")
	@(echo "${OK_COLOR}make vet${NO_COLOR}                      Run go vet on the source code")
	@(echo "${OK_COLOR}make test${NO_COLOR}                     Run the unit tests")
	@(echo "${OK_COLOR}make test-integration${NO_COLOR}         Run the integration tests")
	@(echo "${OK_COLOR}make watch${NO_COLOR}                    Run funzzy watch")
	@(echo "${OK_COLOR}make clean${NO_COLOR}                    Remove the compiled executables")
	@(echo "${OK_COLOR}make deps${NO_COLOR}                     Get the dependencies needed to build the project")

all: deps test build test-integration bump-version

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

tools:
	@(go get github.com/golang/lint)

fmt:
	@(echo "${OK_COLOR}Running fmt ...${NO_COLOR}")
	@([ $$(gofmt -l . | wc -l) != 0 ] && \
	echo "${WARN_COLOR}The following files are not correctly formated:${NO_COLOR}" && \
	echo "${ERROR_COLOR}" && gofmt -l . && \
	echo "${NO_COLOR}"  && exit 1 || exit 0)

vet:
	@(echo "${OK_COLOR}Running vet ...${NO_COLOR}")
	go vet ./...

lint: tools
	@(echo "${OK_COLOR}Running lint ...${NO_COLOR}")
	@(export PATH=$$PATH:$$GOPATH/bin && [ $$(golint ./... | wc -l) != 0 ] && \
	echo "${WARN_COLOR}Lint says the following files are not ok:${NO_COLOR}" && \
	echo "${ERROR_COLOR}" && golint ./... && \
	echo "${NO_COLOR}"  && exit 1 || exit 0)

test: fmt vet lint
	@(go test -v ./...)

test-integration: build
	go test -tags=integration -v ./...

watch:
	funzzy watch

clean:
	rm bin/ergo

deps:
	@go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v ergo | go get
