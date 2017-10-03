
# Ergo [![GoDoc](https://godoc.org/github.com/cristianoliveira/ergo?status.svg)](https://godoc.org/github.com/cristianoliveira/ergo) [![Go Report Card](https://goreportcard.com/badge/github.com/cristianoliveira/ergo)](https://goreportcard.com/report/github.com/cristianoliveira/ergo) [![unix build](https://img.shields.io/travis/cristianoliveira/ergo.svg?label=unix)](https://travis-ci.org/cristianoliveira/ergo) [![win build](https://img.shields.io/appveyor/ci/cristianoliveira/ergo.svg?label=win)](https://ci.appveyor.com/project/cristianoliveira/ergo) 

<p align="left" >
<img src="https://s-media-cache-ak0.pinimg.com/736x/aa/bc/3b/aabc3b2b789f478ffb87ac2f0bdd2d33--ergo-proxy-manga-anime.jpg" width="250" align="center" />
<span>Ergo Proxy - Reverse proxy agent for local domain.</span>

</p>

<p align="center">
  Management of multiple apps running on different ports made easy through local domains.
</p>

## Demo

<img src="https://raw.githubusercontent.com/cristianoliveira/ergo/master/demo.gif" align="center" />

See more [examples](https://github.com/cristianoliveira/ergo/tree/master/examples)

Ergo's goal is to be a simple reverse proxy that follows [Unix's philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) of doing only one thing and doing it well. Simplicity means no magic involved. Just a flexible reverse proxy that extends the well-known `/etc/hosts` declaration.

**Feedback**

This project is still under development but it's ready to use. Feel free to leave your 
feedback and opening issues. Suggestions and contributions are welcome. :)

## Why?

When dealing with multiple apps locally it's really annoying rememberinf which port represents each service and it gets even worse when you have microservices. Ergo solves this by assigning a local domain for each app.

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

I use unix based systems on a daily basis, so I can't test every build on every OS :(

### Using go
```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path. `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` file inside the current directory. It must contain the names and URL of the services following the same format as the `/etc/hosts` (domain+space+url). The main difference is that it also considers the specified port.

Let's start:

### Simplest Setup

**You need to set the `http://127.0.0.1:2000/proxy.pac` configuration on your system network config**

Ergo provides a setup command to generate those config for you. Currently, the supported OS's are:

 - osx
 - linux-gnome
 - windows

(Contributions are welcomed)

```bash
ergo setup <operation-system>
```

If you have trouble using Ergo, please take a look on detailed config session below.

### Showtime

```
echo "ergoproxy http://localhost:3000" > .ergo
ergo run
```
Now you are able to access: `http://ergoproxy.dev`.

Ergo redirects anything that ends with `*.dev` to the configured url.

Simple, right? No magic involved.

You can easily add more services by adding more lines in `.ergo`:
```
echo "otherservice http://localhost:5000" >> .ergo
ergo list
ergo run
```

Restart the Ergo server and access: `http://otherservice.dev`

## Configuration

In order to use Ergo you need to set it up as your machine's proxy. Just point your proxy config to `http://127.0.0.1:2000/proxy.pac`:

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

You can run Ergo with an ephemeral setup, check the scripts inside `/resources` forlder. Those scripts set the proxy config only while `ergo` is running.

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
