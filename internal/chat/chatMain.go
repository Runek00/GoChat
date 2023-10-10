package chat

import (
	"GoChat/api"
	"net/http"
)

var rooom = newRoom()

func InitChat() {
	http.HandleFunc("/chat", rooom.Handler().ServeHTTP)
	http.HandleFunc("/chatroom/", writeMessage)
	// http.HandleFunc("/")
	rooom.run()
}

func writeMessage(w http.ResponseWriter, r *http.Request) {
	usr := api.GetSessionUser(w, r)
	msg := r.FormValue("chat_message")
	msg = usr.Login + ": " + msg
	rooom.SendAll(&msg)
}
