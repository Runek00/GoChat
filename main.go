package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {

	InitDb()
	defer CloseDb()

	InitAuth()

	h1 := func(w http.ResponseWriter, r *http.Request) {
		if !CheckSession(w, r) {
			return
		}
		http.Redirect(w, r, "/hello", http.StatusSeeOther)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/hello", hello)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	if !CheckSession(w, r) {
		return
	}
	tmpl := template.Must(template.ParseFiles("hello.html"))
	tmpl.Execute(w, nil)
}
