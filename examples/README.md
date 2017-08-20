## Examples

Two services `service_one.go` and `service_two.go` running on `8001` and `8002`.


On the first terminal:
```
go run service_one.go
```

On the second terminal:
```
go run service_two.go
```

Finally run ergo: (See installation and configuration instructions)
```
ergo run
```

Then access: `http://serviceone.dev` and `http://servicetwo.dev`
On `./examples/.ergo` are the configured domains.

Simple :)
