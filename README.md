
# Ergo [![GoDoc](https://godoc.org/github.com/cristianoliveira/ergo?status.svg)](https://godoc.org/github.com/cristianoliveira/ergo) [![Go Report Card](https://goreportcard.com/badge/github.com/cristianoliveira/ergo)](https://goreportcard.com/report/github.com/cristianoliveira/ergo) [![unix build](https://img.shields.io/travis/cristianoliveira/ergo.svg?label=unix)](https://travis-ci.org/cristianoliveira/ergo) [![win build](https://img.shields.io/appveyor/ci/cristianoliveira/ergo.svg?label=win)](https://ci.appveyor.com/project/cristianoliveira/ergo) [![codecov](https://codecov.io/gh/cristianoliveira/ergo/branch/master/graph/badge.svg)](https://codecov.io/gh/cristianoliveira/ergo)

<p align="left" >
<img src="https://s-media-cache-ak0.pinimg.com/736x/aa/bc/3b/aabc3b2b789f478ffb87ac2f0bdd2d33--ergo-proxy-manga-anime.jpg" width="250" align="center" />
<span>Ergo Proxy - The reverse proxy agent for local domain management.</span>

</p>

<p align="center">
  The management of multiple apps running over different ports made easy through custom local domains.
</p>

## Demo

<img src="https://raw.githubusercontent.com/cristianoliveira/ergo/master/demo.gif" align="center" />

See more on [examples](https://github.com/cristianoliveira/ergo/tree/master/examples)

Ergo's goal is to be a simple reverse proxy that follows the [Unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) of doing only one thing and doing it well. Simplicity means no magic involved. Just a flexible reverse proxy which extends the well-known `/etc/hosts` declaration.

**Feedback**

This project is constantly undergoing development, however, it's ready to use. Feel free to provide
feedback as well as open issues. All suggestions and contributions are welcome. :)

## Why?

Dealing with multiple apps locally, and having to remember each port representing each microservice is frustrating. I wanted a simple way to assign each service a proper local domain. Ergos solves this problem.

## Installation

### OSX
```
brew tap cristianoliveira/tap
brew install ergo
```

### Linux
```
curl -s https://raw.githubusercontent.com/cristianoliveira/ergo/master/install.sh | sh
```

### Windows
You can find the Windows executables in [release](https://github.com/cristianoliveira/ergo/releases).

***Disclaimer:***
I use Unix-based systems on a daily basis, so I am not able to test each build alone. :(

### Go
```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path: `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` file inside the current directory. It must contain the names and URL of the services following the same format as `/etc/hosts` (`domain`+`space`+`url`). The main difference is it also considers the specified port.

### Simplest Setup

**You need to set the `http://127.0.0.1:2000/proxy.pac` configuration on your system network config.**

Ergo comes with a setup command that can configure it for you. The current systems supported are:

 - osx
 - gnome (tested on Linux and FreeBSD)
 - windows

```bash
ergo setup <operation-system>
```

In case of errors / it doesn't work, please look at the detailed config session below.

### Showtime

```
echo "ergoproxy http://localhost:3000" > .ergo
ergo run
```
Now you should be able to access: `http://ergoproxy.dev`.
Ergo redirects anything ending with `.dev` to the configured url.

Simple, right? No magic involved.

Do you want to add more services? It's easy, just add more lines in `.ergo`:
```
echo "otherservice http://localhost:5000" >> .ergo
ergo list
ergo run
```

Restart the server and access: `http://otherservice.dev`

## Configuration

In order to use Ergo domains you need to set it as a proxy. Set the `http://127.0.0.1:2000/proxy.pac` on:

##### OS X

`Network Preferences > Advanced > Proxies > Automatic Proxy Configuration`

##### Windows

`Settings > Network and Internet > Proxy > Use setup script`

##### Linux

On Ubuntu

`System Settings > Network > Network Proxy > Automatic`

For other distributions, check your network manager and look for proxy configuration. Use browser configuration as an alternative.

### Browser configuration

Browsers can be configured to use a specific proxy. Use this method as an alternative to system-wide configuration.

##### Chrome

Exit Chrome and start it using the following option:

```sh
# Linux
$ google-chrome --proxy-pac-url=http://localhost:2000/proxy.pac

# OS X
$ open -a "Google Chrome" --args --proxy-pac-url=http://localhost:2000/proxy.pac
```

### Ephemeral Setup

As an alternative you can see the scripts inside `/resources` for running an
ephemeral setup. Those scripts set the proxy only while `ergo` is running.

## Run tests

```
  make test
```

## Contributing
 - Fork it!
 - Create your feature branch: `git checkout -b my-new-feature`
 - Commit your changes: `git commit -am 'Add some feature'`
 - Push to the branch: `git push origin my-new-feature`
 - Submit a pull request, they are welcome!
 - Please include unit tests in your pull requests

# License

MIT
