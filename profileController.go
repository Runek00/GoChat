package main

import (
	"html/template"
	"net/http"
)

type ProfileInfo struct {
	User    User
	CanEdit bool
}

func InitProfile() {
	http.HandleFunc("/myProfile", showMyProfile)
	http.HandleFunc("/editProfile", editProfile)
}

func showMyProfile(w http.ResponseWriter, r *http.Request) {
	usr := GetSessionUser(w, r)
	info := ProfileInfo{usr, true}
	tmpl := template.Must(template.ParseFiles("profile.html"))
	tmpl.Execute(w, info)
}

func editProfile(w http.ResponseWriter, r *http.Request) {
	usr := GetSessionUser(w, r)
	switch r.Method {
	case "GET":
		info := ProfileInfo{usr, true}
		tmpl := template.Must(template.ParseFiles("editProfile.html"))
		tmpl.Execute(w, info)
	case "POST":
		usr.Location = r.FormValue("location")
		usr.Info = r.FormValue("info")
		EditUser(usr)
		info := ProfileInfo{usr, true}
		tmpl := template.Must(template.ParseFiles("profile.html"))
		tmpl.Execute(w, info)
	}
}
