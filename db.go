package main

import (
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
	"html/template"
	"errors"
)

type Guide struct {
	Subject, Title string
	Content template.HTML
	Delta, Creator string
}

var db *sql.DB

func DBinit() error {
	var err error
	db, err = sql.Open("sqlite3", "file:./database.sqlite?cache=shared&mode=rwc")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS guides (id INTEGER PRIMARY KEY AUTOINCREMENT, subject TEXT, title TEXT UNIQUE, content TEXT, delta TEXT, creator TEXT)")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE, password TEXT)")
	return err
}
func DBclose() {db.Close()}

// Study guide methods {{{
func DBinsert(subject, title, content, delta, creator string) error {
	_, err := db.Exec("INSERT INTO guides (subject, title, content, delta, creator) VALUES (?, ?, ?, ?, ?)", subject, title, content, delta, creator)
	return err
}
func DBedit(title, content, delta string) error {
	_, err := db.Exec("UPDATE \"guides\" SET content = ?, delta = ? WHERE title = ?", content, delta, title)
	return err
}
func DBget(title string) (Guide, bool, error) {
	rows, err := db.Query("SELECT content, subject, delta, creator FROM guides WHERE title = ?", title)
	if err != nil {
		return Guide{}, false, err
	}
	defer rows.Close()
	if rows.Next() {
		var content, subject, delta, creator string
		err = rows.Scan(&content, &subject, &delta, &creator)
		if err != nil {
			return Guide{}, false, err
		}
		return Guide{subject, title, template.HTML(content), delta, creator}, false, nil
	}
	return Guide{}, true, nil
}
func DBgetAll() ([]Guide, error) {
	rows, err := db.Query("SELECT subject, title, content, delta, creator FROM guides")
	if err != nil {
		return []Guide{}, err
	}
	defer rows.Close()
	returnval := make([]Guide, 0)
	for rows.Next() {
		var c, t, s, d, cr string
		err = rows.Scan(&s, &t, &c, &d, &cr)
		if err != nil {
			return []Guide{}, err
		}
		returnval = append(returnval, Guide{s, t, template.HTML(c), d, cr})
	}
	return returnval, nil
}
// }}}

// User methods: {{{
func DBcreateUser(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashed))
	if err != nil {
		return err
	}
	return nil
}
func DBlogIn(username, password string) (bool, error) {
	rows, err := db.Query("SELECT password FROM users WHERE username = ?", username)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		var hashed string
		err = rows.Scan(&hashed)
		if err != nil {
			return false, err
		}
		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		if err != nil {
			return false, nil
		}
		return true, nil
	}
	return false, errors.New("User not found")
}
// TODO: Combine title and subject into one thing! Never change independently and are a PITA to handle separately
func DBgetOfUser(username string) ([]struct{Title, Subject string}, error) { //returns title/subjects of guides the user has created
	rows, err := db.Query("SELECT title, subject FROM guides WHERE creator = ?", username)
	if err != nil {
		return []struct{Title, Subject string}{}, err
	}
	defer rows.Close()
	returnval := make([]struct{Title, Subject string}, 0)
	for rows.Next() {
		var title, subject string
		err = rows.Scan(&title, &subject)
		if err != nil {
			return []struct{Title, Subject string}{}, err
		}
		returnval = append(returnval, struct{Title, Subject string}{title, subject})
	}
	return returnval, nil
}
// }}}
