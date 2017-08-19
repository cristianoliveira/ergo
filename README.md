# Ergo

The proxy for local microservices managment and development.
The goal of ergo is to be a simple reverse proxy that follow the unix philosophy of doing only one thing.

This project is under development but it is usable. Feel free to give me
feedback and opening issues. Suggestions and contributions are welcome. :)

## Installation

```
go install github.com/cristianoliveira/ergo
```
Make sure you have `$GOPATH/bin` in your path. `export PATH=$PATH:$GOPATH/bin`

## Usage

Ergo looks for a `.ergo` inside the current folder. It must contain the names and
url of the services following the same format as the `/etc/hosts` the main difference
is that Ergo also consider the port specified.

**Ergo runs on `localhost:8080` you have to configure it as your proxy in Network configs of your system**

Let's start:
```
echo "ergoproxy=http://localhost:3000" > .ergo && ergo
```
Now you are able to access: `http://ergoproxy.dev`.
Ergo redirects anything that finish with `.dev` to the configured url.
Simple, no magic involved.

Do you want add more services? So is simple add more lines in `.ergo`:
```
echo "otherservice=http://localhost:5000" >> .ergo
```

Restart the server and access: `http://otherservice.dev`

# License

MIT
