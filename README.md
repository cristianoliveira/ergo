# Ergo

The proxy for local managment of microservices. The goal of ergo is to be
a reverse proxy that follow the unix filosofy of do only one thing.

## Usage

Ergo looks for a `.ergo` in the current folder. It should contain the names and
url of the services following the format: `servicename=url`.
Ergo runs on `localhost:8080` you have to configure it as your proxy.

Let's start:
```
echo "ergoproxy=http://localhost:3000" > .ergo && ergo
```
Now ergo is redirect anything that finish with `ergoproxy.dev` to the configured
url. Simples.

Do you want add more services? So simple add more lines in `.ergo`:
```
echo "otherservice=http://localhost:5000" >> .ergo
```

# License

MIT
