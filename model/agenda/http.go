package agenda

import (
	"log"
	"net/http"
	"text/template"

	"worgenda/app"
)

func ShowMain(w http.ResponseWriter, r *http.Request) {
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("model/agenda/tmpl/agenda.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}
