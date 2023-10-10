package chat

import "net/http"

func InitChat() {
	room := newRoom()
	http.HandleFunc("/chatroom", room.Handler().ServeHTTP)
	room.run()

}
