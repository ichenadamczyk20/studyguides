package main

import (
	"log"
	"net/http"
	"net/url"
	"html/template"
	"io/ioutil"
	"github.com/gorilla/sessions"
	"errors"
)

var sessionStore *sessions.FilesystemStore

var templates = template.Must(template.ParseGlob("./templates/*.html"))
func executeTemplate(w http.ResponseWriter, templ string, content interface{}) {
	err := templates.ExecuteTemplate(w, templ, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	err := DBinit()
	if err != nil {log.Fatal(err)}
	defer DBclose()

	key, err := ioutil.ReadFile("key")
	if err != nil {log.Fatal(err)}
	sessionStore = sessions.NewFilesystemStore("", key)

	http.HandleFunc("/create", create)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/guide/", guide)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Running stuyguides")
}

func home(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		guides, err := DBgetAll()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		executeTemplate(w, "index.html", guides)
	case "/about":
		executeTemplate(w, "about.html", "")
	default:
		http.NotFound(w, r)
	}
}

// Guide methods: {{{
func guide(w http.ResponseWriter, r *http.Request) {
	path, err := url.PathUnescape(r.URL.Path[len("/guide/"):])
	if err != nil {
		http.Error(w, "Internal Server Error" + err.Error(), http.StatusInternalServerError)
		return
	}
	guide, notFound, err := DBget(path)
	if err != nil {
		if notFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Internal Server Error" + err.Error(), http.StatusInternalServerError)
		}
		return
	}
	executeTemplate(w, "guide.html", guide)
}
func edit(w http.ResponseWriter, r *http.Request) {
	path, err := url.PathUnescape(r.URL.Path[len("/edit/"):])
	if err != nil {
		http.Error(w, "Internal Server Error" + err.Error(), http.StatusInternalServerError)
		return
	}
	guide, notFound, err := DBget(path)
	if err != nil {
		if notFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Internal Server Error" + err.Error(), http.StatusInternalServerError)
		}
		return
	}
	switch r.Method {
	case "GET":
		executeTemplate(w, "edit.html", guide)
	case "POST":
		_, err, status := authUser(r) //TODO: changelog
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		if r.FormValue("content") == "" || r.FormValue("delta") == "" {
			http.Error(w, "The updated content of the guide is missing", http.StatusBadRequest)
			return
		}
		err = DBedit(path, r.FormValue("content"), r.FormValue("delta"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/guide/" + path, http.StatusFound)
	}
}
func create(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		executeTemplate(w, "create.html", "")
	case "POST":
		username, err, status := authUser(r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		if r.FormValue("content") == "" || r.FormValue("delta") == "" || r.FormValue("subject") == "" || r.FormValue("title") == "" {
			http.Error(w, "Either the title, the subject or the content are missing", http.StatusBadRequest)
			return
		}
		_, notFound, err := DBget(r.FormValue("title"))
		if err != nil {
			http.Error(w, "Guide already exists; edit it instead", http.StatusBadRequest)
			return
		}
		if !notFound {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = DBinsert(r.FormValue("subject"), r.FormValue("title"), r.FormValue("content"), r.FormValue("delta"), username)
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
// }}}

// User methods: {{{
func authUser(r *http.Request) (string, error, int) { // username, error, status code
	session, err := sessionStore.Get(r, "user")
	if err != nil {
		return "", err, http.StatusInternalServerError //TODO: fix
	}
	interface_username, exists := session.Values["username"]
	interface_password, exists := session.Values["password"]
	if !exists {
		return "", errors.New("You are not logged in!"), http.StatusUnauthorized
	}
	username, password := interface_username.(string), interface_password.(string) //TODO: type switch + error
	loggedIn, err := DBlogIn(username, password)
	if err != nil {
		return "", errors.New("Error logging in: " + err.Error()), http.StatusInternalServerError
	}
	if !loggedIn {
		return "", errors.New("You are not properly logged in!"), http.StatusUnauthorized
	}
	return username, nil, 0
}
// }}}
