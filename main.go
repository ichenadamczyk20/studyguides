package main

import (
	"log"
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	err := DBinit()
	if err != nil {log.Fatal(err)}
	defer DBclose()
	http.HandleFunc("/create", create)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/guide/", guide)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Running stuyguides")
}
func guide(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/guide/"):]
	content, notFound, err := DBget(path)
	if err != nil {
		if notFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Internal Server Error" + err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "You requested the guide: %s\nContent: %s\n", path, content)
}
func edit(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/edit/"):]
	fmt.Fprintf(w, "You are going to edit the guide: %s", path)
}
func create(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "You are creating a guide")
	case "POST":
		if r.FormValue("content") == "" || r.FormValue("title") == "" {
			http.Error(w, "Either the title or the content are missing", http.StatusBadRequest)
			return
		}
		_, notFound, err := DBget(r.FormValue("title"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !notFound {
			http.Error(w, "Guide already exists; edit it instead", http.StatusBadRequest)
			return
		}
		err = DBinsert(r.FormValue("title"), r.FormValue("content"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/guide/" + url.PathEscape(r.FormValue("title")), http.StatusFound)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "Welcome to the stuygides website")
	case "/about":
		fmt.Fprint(w, "Welcome to the about page!")
	default:
		http.NotFound(w, r)
	}
}

