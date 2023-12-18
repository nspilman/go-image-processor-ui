package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

var FS embed.FS

func main() {
	component := hello("Julian")

	http.Handle("/", templ.Handler(component))
	http.HandleFunc("/output.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/dist/output.css")
	})
	http.Handle("/pioneer", templ.Handler(pioneer("hello")))

	fmt.Println("Listening on :4000")
	http.ListenAndServe(":4000", nil)
}
