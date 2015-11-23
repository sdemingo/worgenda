package main

import (
	"appengine"
	"html/template"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	config, err := GetDropboxConfig()
	if err != nil {
		log.Println(err)
		return
	}
	fileContent, err := ReadFile(config, config.Files[0])

	notes := Parse(fileContent)

	a := NewAgenda()
	a.InsertNotes(notes)
	a.Build()

	tmpl := template.Must(template.ParseFiles("tmpl/agenda.html"))
	if err := tmpl.Execute(w, a); err != nil {
		log.Println(err)
		return
	}
}
