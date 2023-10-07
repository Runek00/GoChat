package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id       int
	login    string
	password string
	regDate  time.Time
	from     string
	info     string
	active   bool
}

func CheckUser(login string, password string) bool {
	hash, ok := getPassword(login)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return ok && err == nil
}

func GetUserByLogin(login string) (User, bool) {
	rows, err := Db().Query("select * from user where login = ?", login)
	if err != nil {
		log.Fatal(err)
	}
	return parseRowsToUser(rows)
}

func GetUser(id int) (User, bool) {
	rows, err := Db().Query("select * from user where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	return parseRowsToUser(rows)
}

func parseRowsToUser(rows *sql.Rows) (User, bool) {
	usr := User{}
	if rows.Next() {
		var regdate int64
		err := rows.Scan(&usr.id, &usr.login, &usr.password, &regdate, &usr.from, &usr.info, &usr.active)
		if err != nil {
			log.Fatal(err)
		}
		usr.regDate = time.Unix(regdate, 0)
		return usr, true
	} else {
		return usr, false
	}
}

func getPassword(login string) (string, bool) {
	rows, err := Db().Query("select password from user where login = ?", login)
	if err != nil {
		log.Fatal(err)
	}
	var password string
	if rows.Next() {
		err = rows.Scan(&password)
		if err != nil {
			log.Fatal(err)
		}
		return password, true
	} else {
		return password, false
	}
}

func AddUser(login string, password string) {
	hash, _ := hashPassword(password)
	_, err := Db().Exec("insert into user(login, password, regDate) values (?, ?, ?)", login, hash, time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
