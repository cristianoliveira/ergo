package main

import (
	"net/http"

	"github.com/cristianoliveira/ergo/proxy"
)

func main() {
	config := proxy.LoadConfig("./.ergo")
	proxy := proxy.NewErgoProxy(config)
	http.Handle("/", proxy)
	http.ListenAndServe(":2000", nil)
}
