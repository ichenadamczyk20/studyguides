package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"errors"
)

var db *sql.DB

func DBinit() error {
	var err error
	db, err = sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS guides (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT UNIQUE, content TEXT)")
	return err
}
func DBinsert(title, content string) error {
	_, err := db.Exec("INSERT INTO guides (title, content) VALUES (?, ?)", title, content)
	return err
}
func DBget(title string) (string, bool, error) {
	rows, err := db.Query("SELECT content FROM guides WHERE title = ?", title)
	if err != nil {
		return "", false, err
	}
	if rows.Next() {
		var content string
		err = rows.Scan(&content)
		if err != nil {
			return "", false, err
		}
		return content, false, nil
	}
	return "", true, errors.New("No guide found")
}
func DBclose() {db.Close()}

