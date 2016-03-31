package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var templates map[string]*template.Template

func init() {
	// big thanks to https://elithrar.github.io/article/approximating-html-template-inheritance/
	templates = make(map[string]*template.Template)
	homeTemp := template.Must(template.ParseFiles("views/index.html", "views/layout.html"))
	templates["index"] = homeTemp
	errorTemp := template.Must(template.ParseFiles("views/error.html", "views/layout.html"))
	templates["error"] = errorTemp
}

type page struct {
	Title string
}

type errorPage struct {
	Message string
	Error   pageError
}

type pageError struct {
	Status int
	Stack  string
}

func homepage(w http.ResponseWriter, r *http.Request) {
	p := &page{
		Title: "Express",
	}
	err := templates["index"].ExecuteTemplate(w, "layout", p)
	if err != nil {
		log.Fatal(err)
	}
}

func notFoundError(w http.ResponseWriter, r *http.Request) {
	ep := &errorPage{
		Message: "Not Found",
		Error: pageError{
			Status: 404,
			Stack:  fmt.Sprintf("Error, url not found: %v", r.URL),
		},
	}
	templates["error"].ExecuteTemplate(w, "layout", ep)
}

func main() {
	router := mux.NewRouter()
	dir := http.Dir("./public/stylesheets/")
	handler := http.StripPrefix("/stylesheets/", http.FileServer(dir))
	http.Handle("/stylesheets/", handler)

	router.HandleFunc("/", homepage)
	router.NotFoundHandler = http.HandlerFunc(notFoundError)

	http.Handle("/", router)
	http.ListenAndServe(":3000", nil)
}
