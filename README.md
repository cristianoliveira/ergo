
# Ergo [![GoDoc](https://godoc.org/github.com/cristianoliveira/ergo?status.svg)](https://godoc.org/github.com/cristianoliveira/ergo) [![Go Report Card](https://goreportcard.com/badge/github.com/cristianoliveira/ergo)](https://goreportcard.com/report/github.com/cristianoliveira/ergo) [![unix build](https://img.shields.io/travis/cristianoliveira/ergo.svg?label=unix)](https://travis-ci.org/cristianoliveira/ergo) [![win build](https://img.shields.io/appveyor/ci/cristianoliveira/ergo.svg?label=win)](https://ci.appveyor.com/project/cristianoliveira/ergo) 

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

The Ergo's goal is to be a simple reverse proxy that follows the [Unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) of doing only one thing and do it well. Simplicity means no magic involved. Just a flexible reverse proxy that extends the well-known `/etc/hosts` declaration.

**Feedback**

This project is under development but it's ready to use. Feel free to give me
feedback and opening issues. Suggestions and contributions are welcome. :)

## Why?

When dealing with multiple apps locally it's really annoying having to remember each port that represents each service and it gets even worse when you have microservices. So I wanted a simple way to give each app a proper local domain. Ergo comes to solve this simple problem.

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
For windows you can find the executable in each [release](https://github.com/cristianoliveira/ergo/releases).

***Disclaimer:***
I only use unix based systems on my daily basis, so I can't test each build :(

### Using go
```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path. `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` file inside the current folder. It must contain the names and URL of the services following the same format as the `/etc/hosts` (domain+space+url) the main difference is that it also considers the specified port.

Let's start:

### Simplest Setup

**You need to set the `http://127.0.0.1:2000/proxy.pac` configuration on your system network config**

Ergo comes with a setup command that can configure that for you. The current systems supported are:

 - osx
 - linux-gnome
 - windows

(Contributions are welcomed)

```bash
ergo setup <operation-system>
```

In case of errors or if it doesn't work please take a look on detailed config session below.

### Showtime

```
echo "ergoproxy http://localhost:3000" > .ergo
ergo run
```
Now you are able to access: `http://ergoproxy.dev`.
Ergo redirects anything that finishes with `.dev` to the configured url.

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

## Testing 

Running tests:
```
  make test
```

## Contributing
 - Fork it!
 - Create your feature branch: `git checkout -b my-new-feature`
 - Commit your changes: `git commit -am 'Add some feature'`
 - Push to the branch: `git push origin my-new-feature`
 - Submit a pull request

Pull Requests are welcome!

**Pull Request should have unit tests**

# License

MIT
