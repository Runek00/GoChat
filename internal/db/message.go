package db

import (
	"fmt"
	"log"
	"time"
)

type MsgStats struct {
	Logins string
	Counts string
}

type rawStats struct {
	logins []string
	counts []int
}

var statsCache rawStats
var cacheValid = false

func AddMessage(msg string, userId int) {
	_, err := Db().Exec("insert into message(message, user_id, message_date)values (?, ?, ?)", msg, userId, time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
	cacheValid = false
}

func GetStats(login string) MsgStats {
	stats := getUnformattedStats()
	return formatStats(stats, login)
}

func formatStats(rawStats rawStats, myLogin string) MsgStats {
	stats := MsgStats{}
	for i := 0; i < len(rawStats.counts); i++ {
		login := rawStats.logins[i]
		count := rawStats.counts[i]
		if len(stats.Counts) == 0 {
			stats.Counts = fmt.Sprint(count)
			stats.Logins = "\"" + login + "\""
			continue
		}
		if login == myLogin {
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

func getUnformattedStats() rawStats {
	if cacheValid {
		return statsCache
	}
	result, err := Db().Query("select u.login, count(*) cnt from user u join message m on m.user_id = u.id group by u.login order by cnt desc;")
	if err != nil {
		log.Println(err)
		return rawStats{}
	}

	stats := rawStats{}
	for result.Next() {
		var login string
		var count int
		result.Scan(&login, &count)
		stats.logins = append(stats.logins, login)
		stats.counts = append(stats.counts, count)
	}
	statsCache = stats
	cacheValid = true
	return stats
}
