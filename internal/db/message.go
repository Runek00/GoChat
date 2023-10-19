package db

import (
	"fmt"
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

type MsgStats struct {
	Logins string
	Counts string
}

func newStats() MsgStats {
	return MsgStats{
		Logins: "",
		Counts: "",
	}
}

func GetStats(userId int) MsgStats {
	result, err := Db().Query("select u.id, u.login, count(*) cnt from user u join message m on m.user_id = u.id group by u.id order by cnt desc;")
	if err != nil {
		log.Println(err)
		return newStats()
	}
	stats := newStats()
	for result.Next() {
		var login string
		var count int
		var id int
		result.Scan(&id, &login, &count)
		if len(stats.Counts) == 0 {
			stats.Counts = fmt.Sprint(count)
			stats.Logins = "\"" + login + "\""
			continue
		}
		if id == userId {
			stats.Logins = "\"" + login + "\", " + stats.Logins
			stats.Counts = fmt.Sprint(count) + ", " + stats.Counts
		} else {
			stats.Logins = stats.Logins + ", \"" + login + "\""
			stats.Counts = stats.Counts + ", " + fmt.Sprint(count)
		}
	}
	stats.Counts = "[" + stats.Counts + "]"
	stats.Logins = "[" + stats.Logins + "]"
	return stats

}
