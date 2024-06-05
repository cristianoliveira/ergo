
# Ergo [![GoDoc](https://godoc.org/github.com/cristianoliveira/ergo?status.svg)](https://godoc.org/github.com/cristianoliveira/ergo) [![Go Report Card](https://goreportcard.com/badge/github.com/cristianoliveira/ergo)](https://goreportcard.com/report/github.com/cristianoliveira/ergo) [![unix](https://github.com/cristianoliveira/ergo/actions/workflows/on-push.yml/badge.svg)](https://github.com/cristianoliveira/ergo/actions/workflows/on-push.yml) [![win build](https://img.shields.io/appveyor/ci/cristianoliveira/ergo.svg?label=win)](https://ci.appveyor.com/project/cristianoliveira/ergo) [![codecov](https://codecov.io/gh/cristianoliveira/ergo/branch/master/graph/badge.svg)](https://codecov.io/gh/cristianoliveira/ergo)

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

## TL;DR;
```bash
python3 -m http.server 8800 & # launch a web server listening on 8080
echo "http://localhost:8800 mylocalsite" > .ergo # adds a service "mylocalsite" to ergo
ergo local & # may need sudo since it binds to port 80
curl http://mylocalsite.localhost
```

## Summary
* [Philosophy](#philosophy)
* [Installation](#installation)
  - [osx](#osx)
  - [linux](#linux)
  - [windows](#windows)
  - [nix](#nix)
  - [from source](#go)
* [Usage](#usage)
* [Configuration](#configuration)
* [Testing](#run-tests)
* [Contributing](#contributing)

### Philosophy

Ergo's goal is to be a simple reverse proxy that follows the [Unix philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) of doing only one thing and doing it well. Simplicity means no magic involved. Just a flexible reverse proxy which extends the well-known `/etc/hosts` declaration.

**Feedback**

This project is constantly undergoing development, however, it's ready to use. Feel free to provide
feedback as well as open issues. All suggestions and contributions are welcome. :)

For help and feedback you can find us at #ergo-proxy channel on https://gopher.slack.com

## Why?

Dealing with multiple apps locally, and having to remember each port representing each microservice is frustrating. I wanted a simple way to assign each service a proper local domain. Ergo solves this problem.

## Installation

**Important** These are the only official ways to install ergo.

### OSX
```
brew tap cristianoliveira/tap
brew install ergo
```

### Linux
To install the latest official version
```
curl -s https://raw.githubusercontent.com/cristianoliveira/ergo/master/install.sh | sh
```

Or to install a specific version
```
curl -s https://raw.githubusercontent.com/cristianoliveira/ergo/master/install.sh v0.2.5 | sh
```

### Windows

From powershell run:

```
Invoke-WebRequest https://raw.githubusercontent.com/cristianoliveira/ergo/master/install.ps1 -out ./install.ps1; ./install.ps1
```

_You can also find the Windows executables in [release](https://github.com/cristianoliveira/ergo/releases)._

***Disclaimer:***
I use Unix-based systems on a daily basis, so I am not able to test each build alone. :(

### Nix

Create a `flake.nix` with the following content:

```nix
{
  description = "My ergo nix configuration";

  inputs = { 
    nixpkgs.url = "github:NixOS/nixpkgs";
    mypkgs = {
      url = "github:cristianoliveira/nixpkgs";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, mypkgs, ... }:
  { 
    nixosConfigurations.nixos = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        ({ config, pkgs, ... }: { 
         # Injects mypkgs into nixpkgs as custom
         # and then can be referenced as `pkgs.custom.ergo`
         nixpkgs.overlays = [ 
            (final: prev: { custom = import mypkgs { inherit pkgs; }; })
          ];
        })

        # Exemplo of installing a package from mypkgs
        ({ config, pkgs, ... }: {
          environment.systemPackages = [
            pkgs.custom.ergo
          ];
        })
      ];
    };
  };
}
```

Then run `nix develop` to enter the shell with `ergo` available.

### Go
```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path: `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` file inside the current directory. It must contain the names and URL of the services following the same format as `/etc/hosts` (`domain`+`space`+`url`). The main difference is it also considers the specified port.

### Subdomains for localhost
 
Run `ergo local` it'll attempt to bind to `localhost:80` and listen for requests to your services as "subdmains" eg. `http://serviceone.localhost` and `http://servicetwo.localhost`. (Check [examples](https://github.com/cristianoliveira/ergo/tree/master/examples) for more)

**Note:** It __may__ requires sudo to bind to port 80.

You can give it a different port by `ergo local -p <port>` and access it through `http://serviceone.localhost:<port>`.

You can also add a different loopback in `/etc/hosts` like `echo '127.0.0.1 localapp' >> /etc/hosts` and run `ergo local -domain localapp` to access your services through `http://serviceone.localapp` and `http://servicetwo.localapp`.

### Setting up as a webproxy

**You need to set the `http://127.0.0.1:2000/proxy.pac` configuration on your system network config.**

Ergo comes with a setup command that can configure it for you. The current systems supported are:

 - osx
 - linux-gnome
 - windows

```bash
ergo setup <operation-system>
```

In case of errors / it doesn't work, please look at the detailed config session below.

### Adding Services and Running

#### OS X / Linux
```
echo "ergoproxy http://localhost:3000" > .ergo
ergo run
```

Now you should be able to access: `http://ergoproxy.dev`.
Ergo redirects anything ending with `.dev` to the configured URL.

#### Windows
You should not use the default `.dev` domain, we suggest `.test` instead (see [#58](https://github.com/cristianoliveira/ergo/issues/58)) unless your service supports https out of the box and you have already a certificate
```
set ERGO_DOMAIN=.test
echo "ergoproxy  http://localhost:3000" > .ergo
ergo list # you shouldn't see any quotas in the output
ergo run
```
Now you should be able to access: `http://ergoproxy.test`.
Ergo redirects anything ending with `.test` to the configured URL.

Simple, right? No magic involved.

Do you want to add more services? It's easy, just add more lines in `.ergo`:
```
echo "otherservice http://localhost:5000" >> .ergo
ergo list
ergo run
```

Restart the ergo server and access: `http://otherservice.dev`

`ergo add otherservice http://localhost:5000` is a shorthand for appending lines to `./.ergo`

### Ergo's configuration

Ergo accepts different configurations like run in different `port` (default: 2000) and change `domain` (default: dev). You can find all this configs on ergo's help running `ergo -h`.

## Configuration

In order to use Ergo domains you need to set it as a proxy. Set the `http://127.0.0.1:2000/proxy.pac` on:

### Networking Web Proxy

#### OS X

`Network Preferences > Advanced > Proxies > Automatic Proxy Configuration`

#### Windows

`Settings > Network and Internet > Proxy > Use setup script`

#### Linux

On Ubuntu

`System Settings > Network > Network Proxy > Automatic`

For other distributions, check your network manager and look for proxy configuration. Use browser configuration as an alternative.

### Browser configuration

Browsers can be configured to use a specific proxy. Use this method as an alternative to system-wide configuration.

Keep in mind that if you requested the site before setting the proxy properly, you have to reset the cache of the browser or change the name of the service. In `incognito` windows cache is disabled by default, so you can use them if you don't wish to delete the cache

Also you should not use the default `.dev` domain, we suggest `.test` instead (see [#58](https://github.com/cristianoliveira/ergo/issues/58)) unless your service supports https out of the box and you have already a certificate

#### Chrome

Exit Chrome and start it using the following option:

```sh
# Linux
$ google-chrome --proxy-pac-url=http://localhost:2000/proxy.pac

# OS X
$ open -a "Google Chrome" --args --proxy-pac-url=http://localhost:2000/proxy.pac
```

#### Firefox

##### through menus and mouse
1. Click the hamburger button otherwise click on "Edit" Menu
1. then "Preferences"
1. then "Settings" button at the bottom of the page ("General" active in sidebar) with title "Network Settings"
1. check `Automatic Proxy configuration URL` and enter value `http://localhost:2000/proxy.pac` below
1. hit "ok"


##### from about:config
`network.proxy.autoconfig_url` -> `http://localhost:2000/proxy.pac`


### Using on terminal

In order to use ergo as your web proxy on terminal you must set the `http_proxy` variable. (Only for linux/osx)

```sh
export http_proxy="http://localhost:2000"
```

### Ephemeral Setup

As an alternative you can see the scripts inside `/resources` for running an
ephemeral setup. Those scripts set the proxy only while `ergo` is running.

## Contributing
 - Fork it!
 - Create your feature branch: `git checkout -b my-new-feature`
 - Commit your changes: `git commit -am 'Add some feature'`
 - Push to the branch: `git push origin my-new-feature`
 - Submit a pull request, they are welcome!
 - Please include unit tests in your pull requests

## Development

Minimal required golang version `go1.22`.
We recommend using [GVM](https://github.com/moovweb/gvm) for managing
your go versions.

Then simply run:
```sh
gvm use $(cat .gvmrc)
```

### Building

```sh
  make all
```

## Testing

  ```sh
  make test
  make test-integration # Requires admin permission so use it carefully.
```

# License

MIT
