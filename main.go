package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	// component := hello("John")
	component := hello("Julian")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :4000")
	http.ListenAndServe(":4000", nil)
}
