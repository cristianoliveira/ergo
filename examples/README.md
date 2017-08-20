## Examples

There are two services `serviceone` and `servicetwo` running on `8001` and `8002`.

On the first terminal:
```bash
go run serviceone/main.go
```

On the second terminal:
```bash
go run servicetwo/main.go
```

Finally run ergo: (See installation and configuration instructions)
```bash
ergo list # To see the configurations
ergo run
```

Then access: `http://serviceone.dev` and `http://servicetwo.dev`
On `./examples/.ergo` are the configured domains.

Simple :)
