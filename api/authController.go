package api

import (
	"GoChat/internal/db"
	"html/template"
	"net/http"
	"os"

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
		tmpl := template.Must(template.ParseFiles("web/templates/register.html"))
		tmpl.Execute(w, nil)
	case "POST":
		_, exists := db.GetUserByLogin(r.FormValue("login"))
		if exists {
			tmpl := template.Must(template.ParseFiles("web/templates/register.html"))
			tmpl.Execute(w, "User with that login already exists")
			return
		}
		db.AddUser(r.FormValue("login"), r.FormValue("password"))
		http.Redirect(w, r, "/hello", http.StatusSeeOther)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
		tmpl.Execute(w, nil)
	case "POST":
		if db.CheckUser(r.FormValue("login"), r.FormValue("password")) {
			createSession(w, r, r.FormValue("login"))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
			tmpl.Execute(w, "Wrong login or password")
		}
	}
}

func createSession(w http.ResponseWriter, r *http.Request, login string) *sessions.Session {
	user, exists := db.GetUserByLogin(login)
	if !exists {
		http.Error(w, "No such user", http.StatusInternalServerError)
		return nil
	}
	closeSessionIfExists(w, r)
	session, err := Store.New(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	session.Values["login"] = login
	session.Values["userId"] = user.Id
	session.Save(r, w)
	return session
}

func closeSessionIfExists(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	if session != nil {
		session.Options.MaxAge = -1
		Store.Save(r, w, session)
	}
}

func CheckSession(w http.ResponseWriter, r *http.Request) bool {
	session, _ := Store.Get(r, "session")
	_, ok := session.Values["userId"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	return ok
}

func GetSessionUser(w http.ResponseWriter, r *http.Request) db.User {
	session, _ := Store.Get(r, "session")
	id, ok := session.Values["userId"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return db.User{}
	}
	usr, ok := db.GetUser(id.(int))
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return db.User{}
	}
	return usr
}
