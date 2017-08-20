package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from service one!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Running service on localhost:8001")
	http.ListenAndServe(":8001", nil)
}
