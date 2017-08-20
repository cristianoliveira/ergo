## Examples

Two services `service_one.go` and `services_two.go` running on `8001` and `8002`.


in one terminal do:
```
go run service_one.go
```

in second terminal do:
```
go run service_two.go
```

Finally run ergo: (See installation and configuration instructions)
```
ergo run
```

Then access: `http://serviceone.dev` and `http://servicetwo.dev`

Simple :)
