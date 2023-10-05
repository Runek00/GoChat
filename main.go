package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func main() {

	InitDb()
	defer CloseDb()

	initStore()

	h1 := func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(w, r) {
			return
		}
		http.Redirect(w, r, "/hello", http.StatusSeeOther)
	}
	http.HandleFunc("/", h1)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/hello", hello)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initStore() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
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
			createSession(w, r, r.FormValue("login"))
			http.Redirect(w, r, "/hello", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	tmpl := template.Must(template.ParseFiles("hello.html"))
	tmpl.Execute(w, nil)
}

func createSession(w http.ResponseWriter, r *http.Request, login string) *sessions.Session {
	session, err := store.New(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	session.Values["login"] = login
	session.Save(r, w)
	return session
}

func checkSession(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["login"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	return ok
}
