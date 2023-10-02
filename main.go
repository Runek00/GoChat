package main

import (
	"log"
	"net/http"
	"text/template"
)

var Texts = make([]string, 0)

func main() {

	InitDb()

	defer CloseDb()

	h1 := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/hello", hello)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("register.html"))
		tmpl.Execute(w, nil)
	case "POST":
		_, exists := GetUser(r.FormValue("login"))
		if exists {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
		AddUser(r.FormValue("login"), r.FormValue("password"))
		http.Redirect(w, r, "/hello", http.StatusSeeOther)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("login.html"))
		tmpl.Execute(w, nil)
	case "POST":
		if CheckUser(r.FormValue("login"), r.FormValue("password")) {
			http.Redirect(w, r, "/hello", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("hello.html"))
	tmpl.Execute(w, nil)
}
