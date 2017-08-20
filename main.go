package main

import (
	"net/http"

	"github.com/cristianoliveira/ergo/proxy"
)

func main() {
	config := proxy.LoodConfig()
	proxy := proxy.NewErgoProxy(config)
	http.Handle("/", proxy)
	http.ListenAndServe(":2000", nil)
}
