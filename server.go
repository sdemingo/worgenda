package main

import (
	"html/template"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	notes, err := ParseFile("test.org")
	if err != nil {
		log.Println(err)
		return
	}

	a := NewAgenda()
	a.InsertNotes(notes)
	a.Build()

	tmpl := template.Must(template.ParseFiles("tmpl/agenda.html"))
	if err := tmpl.Execute(w, a); err != nil {
		log.Println(err)
		return
	}
}
