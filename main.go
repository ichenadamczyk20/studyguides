package main

import (
	"log"
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/create", create)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/guide/", guide)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Running stuyguides")
}
func guide(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/guide/"):]
	fmt.Fprintf(w, "You requested the guide: %s", path)
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
		//add to database here
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


