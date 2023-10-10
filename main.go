package main

import (
	"GoChat/api"
	"GoChat/internal/chat"
	"GoChat/internal/db"
)

func main() {

	db.InitDb()
	defer db.CloseDb()
	api.InitAll()
	chat.InitChat()
}
