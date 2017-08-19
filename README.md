# Ergo

The proxy for local microservices managment and development.
The goal of ergo is to be a simple reverse proxy that follow the unix philosophy of doing only one thing.

## Installation

```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path. `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` inside the current folder. It should contain the names and
url of the services following the same format as the `/etc/host` the main difference
is that Ergo also consider the port specified.

** Ergo runs on `localhost:8080` you have to configure it as your proxy. **

Let's start:
```
echo "ergoproxy=http://localhost:3000" > .ergo && ergo
```
Now ergo is redirect anything that finish with `ergoproxy.dev` to the configured
url. Simples.

Do you want add more services? So is simple add more lines in `.ergo`:
```
echo "otherservice=http://localhost:5000" >> .ergo
```

# License

MIT
