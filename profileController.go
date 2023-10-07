package main

import (
	"net/http"
	"text/template"
)

type profileInfo struct {
	user    User
	canEdit bool
}

func InitProfile() {
	http.HandleFunc("/myProfile", showMyProfile)
}

func showMyProfile(w http.ResponseWriter, r *http.Request) {
	usr := GetSessionUser(w, r)
	info := profileInfo{usr, true}
	tmpl := template.Must(template.ParseFiles("profile.html"))
	tmpl.Execute(w, info)
}
