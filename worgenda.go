package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var agenda *Agenda

func main() {

	drpConfig, err := GetDropboxConfig()
	if err != nil {
		log.Printf("%v", err)
		return
	}

	agenda = NewAgenda()

	go loadEvents(drpConfig, agenda)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/events", eventHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func loadEvents(config *DropboxConfig, a *Agenda) {

	fileContent, err := ReadFile(config, "prueba.org")
	if err != nil {
		log.Panic("%v", err)
		return
	}
	notes := Parse(fileContent)
	a.InsertNotes(notes)
	a.Build()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("tmpl/agenda.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}

func eventHandler(w http.ResponseWriter, r *http.Request) {

	json, err := json.Marshal(agenda.AllNotes)

	if err != nil {
		log.Printf("%v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
