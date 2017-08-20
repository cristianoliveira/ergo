
# Ergo [![Build Status](https://travis-ci.org/cristianoliveira/apitogo.svg?branch=master)](https://travis-ci.org/cristianoliveira/apitogo)

<p align="left" >
<img src="https://s-media-cache-ak0.pinimg.com/736x/aa/bc/3b/aabc3b2b789f478ffb87ac2f0bdd2d33--ergo-proxy-manga-anime.jpg" width="250" align="center" />
<span>Ergo Proxy - The local proxy agent for multiple services development</span>
</p>

The managment of multiple apps running over diferent ports made easy.

The Ergo's goal is to be a simple reverse proxy that follows the unix philosophy of doing only one thing.

Simple means no magic involved.

**Disclaimer**

This project is under development but it's already usable. Feel free to give me
feedback and opening issues. Suggestions and contributions are welcome. :)

## Installation

```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path. `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` file inside the current folder. It must contain the names and
url of the services following the same format as the `/etc/hosts` the main difference
is that Ergo also consider the port specified.

**Ergo runs on `127.0.0.1:2000` you have to configure it as your proxy in Network configs of your system**

Let's start:
```
echo "ergoproxy http://localhost:3000" > .ergo && ergo
```
Now you are able to access: `http://ergoproxy.dev`.
Ergo redirects anything that finish with `.dev` to the configured url.
Simple, no magic involved.

Do you want add more services? So is simple, just add more lines in `.ergo`:
```
echo "otherservice http://localhost:5000" >> .ergo
```

Restart the server and access: `http://otherservice.dev`

# License

MIT
