package main

import (
	"fmt"
	"github.com/cristianoliveira/ergo/proxy"
	"net/http"
)

func main() {
	config := proxy.LoadConfig("./.ergo")

	http.HandleFunc("/proxy.pac", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./resources/proxy.pac")
	})

	proxy := proxy.NewErgoProxy(config)
	http.Handle("/", proxy)

	fmt.Println("Ergo Proxy listening on port: 2000")
	http.ListenAndServe(":2000", nil)
}
