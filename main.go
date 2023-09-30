package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	fmt.Println("hello world")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", h1)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
