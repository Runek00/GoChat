package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	login    string
	password string
}

var db *sql.DB

func InitDb() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Query("select 1 from user")
	if err != nil {
		db.Exec("create table if not exists user(id INTEGER PRIMARY KEY AUTOINCREMENT, login text not null, password text not null);")
		db.Exec("create unique index if not exists user_login_IDX on user (login);")
	}
}

func CloseDb() {
	db.Close()
}

func CheckUser(login string, password string) bool {
	usr, ok := GetUser(login)
	return ok && usr.password == password
}

func GetUser(login string) (User, bool) {
	rows, err := db.Query("select login, password from user where login = ?", login)
	if err != nil {
		log.Fatal(err)
	}
	usr := User{}
	if rows.Next() {
		err = rows.Scan(&usr.login, &usr.password)
		if err != nil {
			log.Fatal(err)
		}
		return usr, true
	} else {
		return usr, false
	}
}

func AddUser(login string, password string) {
	_, err := db.Exec("insert into user(login, password) values (?, ?)", login, password)
	if err != nil {
		log.Fatal(err)
	}
}
