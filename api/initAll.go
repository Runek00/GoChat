package api

import (
	"GoChat/internal/db"
	"html/template"
	"log"
	"net/http"
	"strings"
	txt "text/template"
)

func InitAll() {
	InitAuth()
	InitProfile()
	Init()
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	if !CheckSession(w, r) {
		return
	}
	http.Redirect(w, r, "/hello", http.StatusSeeOther)
}

func Init() {
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/img/", img)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/charts", chartHandler)
	http.HandleFunc("/chart_script", chartScript)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/favicon.ico")
}

func img(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/img/")
	http.ServeFile(w, r, "web/static/img/"+fileName)
}

func hello(w http.ResponseWriter, r *http.Request) {
	if !CheckSession(w, r) {
		return
	}
	tmpl := template.Must(template.ParseFiles("web/templates/hello.html"))
	tmpl.Execute(w, nil)
}

func chartHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/charts.html"))
	tmpl.Execute(w, nil)
}

func chartScript(w http.ResponseWriter, r *http.Request) {
	usr := GetSessionUser(w, r)
	stats := db.GetStats(usr.Id)
	tmpl := txt.Must(txt.ParseFiles("web/templates/chartUpdate.html"))
	tmpl.Execute(w, stats)
}
