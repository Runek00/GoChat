package main

import (
	"log"
	"net/http"
	"text/template"
)

var Texts = make([]string, 0)

func main() {

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, Texts)
	}
	http.HandleFunc("/", h1)

	h2 := func(w http.ResponseWriter, r *http.Request) {
		txt := r.PostFormValue("in1")
		Texts = append(Texts, txt)
		tmpl := template.Must(template.New("t").Parse("<p> {{ . }} </p>"))
		tmpl.Execute(w, txt)
	}
	http.HandleFunc("/add-txt/", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
