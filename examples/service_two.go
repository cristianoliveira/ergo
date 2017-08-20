package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from service two!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Running service on localhost:8002")
	http.ListenAndServe(":8002", nil)
}
