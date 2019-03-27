package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"html/template"
)

type Guide struct {
	Subject, Title string
	Content template.HTML
	Delta string
}

var db *sql.DB

func DBinit() error {
	var err error
	db, err = sql.Open("sqlite3", "file:./database.sqlite?cache=shared&mode=rwc")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS guides (id INTEGER PRIMARY KEY AUTOINCREMENT, subject TEXT, title TEXT UNIQUE, content TEXT, delta TEXT)")
	return err
}
func DBinsert(subject, title, content, delta string) error {
	_, err := db.Exec("INSERT INTO guides (subject, title, content, delta) VALUES (?, ?, ?, ?)", subject, title, content, delta)
	return err
}
func DBedit(title, content, delta string) error {
	_, err := db.Exec("UPDATE \"guides\" SET content = ?, delta = ? WHERE title = ?", content, delta, title)
	return err
}
func DBget(title string) (Guide, bool, error) {
	rows, err := db.Query("SELECT content, subject, delta FROM guides WHERE title = ?", title)
	if err != nil {
		return Guide{}, false, err
	}
	defer rows.Close()
	if rows.Next() {
		var content, subject, delta string
		err = rows.Scan(&content, &subject, &delta)
		if err != nil {
			return Guide{}, false, err
		}
		return Guide{subject, title, template.HTML(content), delta}, false, nil
	}
	return Guide{}, true, nil
}
func DBgetAll() ([]Guide, error) {
	rows, err := db.Query("SELECT subject, title, content, delta FROM guides")
	if err != nil {
		return []Guide{}, err
	}
	defer rows.Close()
	returnval := make([]Guide, 0)
	for rows.Next() {
		var c, t, s, d string
		err = rows.Scan(&s, &t, &c, &d)
		if err != nil {
			return []Guide{}, err
		}
		returnval = append(returnval, Guide{s, t, template.HTML(c), d})
	}
	return returnval, nil
}

func DBclose() {db.Close()}

