package main

import (
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitAuth() {
	initStore()
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
}

func initStore() {
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("register.html"))
		tmpl.Execute(w, nil)
	case "POST":
		_, exists := GetUserByLogin(r.FormValue("login"))
		if exists {
			tmpl := template.Must(template.ParseFiles("register.html"))
			tmpl.Execute(w, "User with that login already exists")
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
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			tmpl := template.Must(template.ParseFiles("login.html"))
			tmpl.Execute(w, "Wrong login or password")
		}
	}
}

func createSession(w http.ResponseWriter, r *http.Request, login string) *sessions.Session {
	user, exists := GetUserByLogin(login)
	if !exists {
		http.Error(w, "No such user", http.StatusInternalServerError)
		return nil
	}
	session, err := Store.New(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	session.Values["login"] = login
	session.Values["userId"] = user.id
	session.Save(r, w)
	return session
}

func CheckSession(w http.ResponseWriter, r *http.Request) bool {
	session, _ := Store.Get(r, "session")
	_, ok := session.Values["userId"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	return ok
}

func GetSessionUser(w http.ResponseWriter, r *http.Request) User {
	session, _ := Store.Get(r, "session")
	id, ok := session.Values["userId"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return User{}
	}
	usr, ok := GetUser(id.(int))
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return User{}
	}
	return usr
}
