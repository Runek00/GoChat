package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Login    string
	password string
	RegDate  time.Time
	Location string
	Info     string
	Active   bool
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
	defer rows.Close()
	return parseRowsToUser(rows)
}

func GetUser(id int) (User, bool) {
	rows, err := Db().Query("select * from user where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return parseRowsToUser(rows)
}

func parseRowsToUser(rows *sql.Rows) (User, bool) {
	usr := User{}
	if rows.Next() {
		var regdate int64
		err := rows.Scan(&usr.Id, &usr.Login, &usr.password, &regdate, &usr.Location, &usr.Info, &usr.Active)
		if err != nil {
			log.Fatal(err)
		}
		usr.RegDate = time.Unix(regdate, 0)
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
	defer rows.Close()
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
	_, err := Db().Exec("insert into user(login, password, regDate, location, info, active) values (?, ?, ?, ?, ?, ?)", login, hash, time.Now().Unix(), "", "", true)
	if err != nil {
		log.Fatal(err)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func EditUser(user User) {
	_, err := Db().Exec("update user set password = ?, location = ?, info=?, active = ? where id = ?", user.password, user.Location, user.Info, user.Active, user.Id)
	if err != nil {
		log.Fatal(err)
	}
}
