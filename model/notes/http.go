package notes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
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

func NewEventForm(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.NotFound(w, r)
		return
	}
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	var contents map[string]interface{}

	r.ParseForm()
	date, err := time.Parse(DATEFORMATPRINT, r.FormValue("date"))
	if err == nil {
		contents = map[string]interface{}{
			"StringDate": date.Format(DATEFORMATFORHTML),
			"Date":       date,
		}
	}

	// Write template
	tmpl := template.Must(template.ParseFiles("model/notes/tmpl/new-event.html"))
	if err := tmpl.Execute(w, contents); err != nil {
		log.Printf("%v", err)
		return
	}
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.NotFound(w, r)
		return
	}
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	r.ParseForm()
	date, err := time.Parse(DATEFORMATPRINT, r.FormValue("date"))
	if err != nil {
		log.Printf("notes: getevent: bad date: %v", err)
		return
	}

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

func GetEvent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.NotFound(w, r)
		return
	}
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	r.ParseForm()
	date, err := time.Parse(DATEFORMATPRINT, r.FormValue("date"))
	if err != nil {
		log.Printf("notes: getevent: bad date: %v", err)
		return
	}

	id64, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		log.Printf("notes: getevent: bad id: %v", err)
		return
	}

	id := int(id64)
	if id < 0 || id >= len(AllNotes) {
		log.Printf("notes: getevent: bad id: %v", err)
		return
	}

	var contents = map[string]interface{}{
		"StringDate": date.Format(DATEFORMATFORHTML),
		"Date":       date,
		"Event":      &AllNotes[id],
	}

	// Write template
	tmpl := template.Must(template.ParseFiles("model/notes/tmpl/event.html"))
	if err := tmpl.Execute(w, contents); err != nil {
		log.Printf("%v", err)
		return
	}
}

func GetMarkDates(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.NotFound(w, r)
		return
	}
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
