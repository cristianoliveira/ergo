package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from service one!")
	})
	fmt.Println("Running service on localhost:8001")
	http.ListenAndServe(":8001", nil)
}
