package notes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
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

func GetEvents(w http.ResponseWriter, r *http.Request) {
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String() // Does a complete copy of the bytes in the buffer.
	date, _ := time.Parse(DATEFORMATPRINT, s)

	//notes := make([]Note, 0)
	dayNotes := NewDayNotes(date)
	for _, note := range AllNotes {
		if note.InDay(date) {
			dayNotes.Add(note)
		}
	}

	sort.Sort(dayNotes)
	var contents = map[string]interface{}{
		"StringDate": date.Format(DATEFORMATFORHTML),
		"Date":       date,
		"Events":     dayNotes.Notes,
	}

	// Write template
	tmpl := template.Must(template.ParseFiles("model/notes/tmpl/day-events.html"))
	if err := tmpl.Execute(w, contents); err != nil {
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
