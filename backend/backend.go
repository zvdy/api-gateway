package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from backend")
	})

	http.ListenAndServe(":5000", nil)
}
