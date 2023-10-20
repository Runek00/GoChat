package db

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func InitDb() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	createDb(db)
}

func createDb(db *sql.DB) {
	files, err := os.ReadDir("sql")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		sql, err := os.ReadFile("sql/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		db.Exec(string(sql))
	}
}

func CloseDb() {
	db.Close()
}

func Db() *sql.DB {
	return db
}
