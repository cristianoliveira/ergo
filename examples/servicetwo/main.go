package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from service two!")
	})
	fmt.Println("Running service on localhost:8002")
	http.ListenAndServe(":8002", nil)
}
