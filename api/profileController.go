package api

import (
	"GoChat/internal/db"
	"html/template"
	"net/http"
)

type ProfileInfo struct {
	User    db.User
	CanEdit bool
}

func InitProfile() {
	http.HandleFunc("/myProfile", showMyProfile)
	http.HandleFunc("/editProfile", editProfile)
}

func showMyProfile(w http.ResponseWriter, r *http.Request) {
	usr := GetSessionUser(w, r)
	info := ProfileInfo{usr, true}
	tmpl := template.Must(template.ParseFiles("web/templates/profile.html"))
	tmpl.Execute(w, info)
}

func editProfile(w http.ResponseWriter, r *http.Request) {
	usr := GetSessionUser(w, r)
	switch r.Method {
	case "GET":
		info := ProfileInfo{usr, true}
		tmpl := template.Must(template.ParseFiles("web/templates/editProfile.html"))
		tmpl.Execute(w, info)
	case "POST":
		usr.Location = r.FormValue("location")
		usr.Info = r.FormValue("info")
		db.EditUser(usr)
		info := ProfileInfo{usr, true}
		tmpl := template.Must(template.ParseFiles("web/templates/profile.html"))
		tmpl.Execute(w, info)
	}
}
