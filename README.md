
# Ergo [![Build Status](https://travis-ci.org/cristianoliveira/apitogo.svg?branch=master)](https://travis-ci.org/cristianoliveira/apitogo)

<p align="left" >
<img src="https://s-media-cache-ak0.pinimg.com/736x/aa/bc/3b/aabc3b2b789f478ffb87ac2f0bdd2d33--ergo-proxy-manga-anime.jpg" width="250" align="center" />
<span>Ergo Proxy - The local proxy agent for multiple services development</span>
</p>

The managment of multiple apps running over diferent ports made easy.

The Ergo's goal is to be a simple reverse proxy that follows the [unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) of doing only one thing and do it well.

Simplicity means no magic involved. Just a flexible reverse proxy.

**Disclaimer**

This project is under development but it's already usable. Feel free to give me
feedback and opening issues. Suggestions and contributions are welcome. :)

## Why?

When dealing with multiple apps locally is really annoyoing have to remember each
port that represent each service. So I wanted a simple way to give a proper local
domain for each app. Ergos comes to solve this simple problem.

It's not aim to be fancy. It solve this problem and nothing else.
Do you want a web interface? You can either try other projects or create it
on top of ergo's interface. That's the magic of the unix philosophy, composition. :D

## Installation

```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path. `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` file inside the current folder. It must contain the names and
url of the services following the same format as the `/etc/hosts` the main difference
is that Ergo also consider the port specified.

**Set the `http://127.0.0.1:2000/proxy.pac` configuration on your system network config (Details below)**

Let's start:
```
echo "ergoproxy http://localhost:3000" > .ergo
ergo run
```
Now you are able to access: `http://ergoproxy.dev`.
Ergo redirects anything that finish with `.dev` to the configured url.
Simple, no magic involved.

Do you want add more services? So is simple, just add more lines in `.ergo`:
```
echo "otherservice http://localhost:5000" >> .ergo
ergo run
```

Restart the server and access: `http://otherservice.dev`

## Configuration

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

# License

MIT
