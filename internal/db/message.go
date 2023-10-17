package db

import (
	"log"
	"time"
)

// type message struct {
// 	msg    string
// 	userId int
// 	date   time.Time
// }

func AddMessage(msg string, userId int) {
	_, err := Db().Exec("insert into message(message, user_id, message_date)values (?, ?, ?)", msg, userId, time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
}
