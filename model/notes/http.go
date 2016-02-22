package notes

import (
	"log"
	"net/http"
	"text/template"

	"worgenda/app"
)

func Main(w http.ResponseWriter, r *http.Request) {
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("model/notes/tmpl/agenda.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}
