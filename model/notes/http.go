package notes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

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

func GetMarkDates(w http.ResponseWriter, r *http.Request) {
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	dates := make(map[time.Time]bool)

	for _, note := range AllNotes {
		for _, stamp := range note.Stamps {
			if ok, _ := dates[stamp]; !ok {
				dates[stamp] = true
			}
		}
	}

	alldates := make([]string, 0, len(dates))
	for date := range dates {
		sdate := date.Format(DATEFORMATPRINT)
		alldates = append(alldates, sdate)
	}

	// Write json
	w.Header().Set("Content-Type", "application/json")
	jbody, err := json.Marshal(alldates)
	if err != nil {
		log.Printf("notes: getmarkdates: %v", err)
		return
	}
	fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))
}
