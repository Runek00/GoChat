package db

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitDb() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Query("select 1 from user")
	if err != nil {
		db.Exec("create table if not exists user(id INTEGER PRIMARY KEY AUTOINCREMENT, login text not null, password text not null, regdate integer, location text, info text, active boolean);")
		db.Exec("create unique index if not exists user_login_IDX on user (login);")
	}
}

func CloseDb() {
	db.Close()
}

func Db() *sql.DB {
	return db
}
