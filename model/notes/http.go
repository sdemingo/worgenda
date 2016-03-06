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

	tmpl := template.Must(template.ParseFiles(app.AppDir + "/model/notes/tmpl/agenda.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}

func NewEventForm(w http.ResponseWriter, r *http.Request) {
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	var contents map[string]interface{}

	r.ParseForm()
	date, err := time.Parse(DATEFORMATGETPARAM, r.FormValue("date"))
	if err == nil {
		contents = map[string]interface{}{
			"StringDate":       date.Format(DATEFORMATFORHTML),
			"SimpleStringDate": date.Format(DATEFORMATGETPARAM),
			"Date":             date,
		}
	}

	// Write template
	tmpl := template.Must(template.ParseFiles(app.AppDir + "/model/notes/tmpl/new-event.html"))
	if err := tmpl.Execute(w, contents); err != nil {
		log.Printf("%v", err)
		return
	}
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		http.NotFound(w, r)
		return
	}

	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	newNote := struct {
		Title string
		Date  string
		Hour  string
		Body  string
	}{"", "", "", ""}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newNote)
	if err != nil {
		log.Printf("notes: addevent: %v", err)
		return
	}

	note := NewNote()
	note.Title = newNote.Title
	note.Source = "worgenda.org"
	note.Body = newNote.Body
	stamp, err := time.Parse(DATEFORMATGETPARAM+" "+HOURFORMATPRINT, newNote.Date+" "+newNote.Hour)
	if err == nil {
		note.Stamps = append(note.Stamps, stamp)
	}

	fmt.Println(note)

	// Write json
	w.Header().Set("Content-Type", "application/json")
	jbody, err := json.Marshal(newNote)
	if err != nil {
		log.Printf("notes: getmarkdates: %v", err)
		return
	}
	fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))
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
	AllNotes.GetNotesFromDate(dayNotes)

	sort.Sort(dayNotes)
	var contents = map[string]interface{}{
		"FormParamDate": date.Format(DATEFORMATGETPARAM),
		"StringDate":    date.Format(DATEFORMATFORHTML),
		"Date":          date,
		"Events":        dayNotes.Notes,
	}

	// Write template
	tmpl := template.Must(template.ParseFiles(app.AppDir + "/model/notes/tmpl/day-events.html"))
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

	note := AllNotes.GetNote(id64)
	if note == nil {
		log.Printf("notes: getevent: bad id: %v", err)
		return
	}

	var contents = map[string]interface{}{
		"StringDate": date.Format(DATEFORMATFORHTML),
		"Date":       date,
		"Event":      note,
	}

	// Write template
	tmpl := template.Must(template.ParseFiles(app.AppDir + "/model/notes/tmpl/event.html"))
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

	alldates := AllNotes.GetBusyDates()

	// Write json
	w.Header().Set("Content-Type", "application/json")
	jbody, err := json.Marshal(alldates)
	if err != nil {
		log.Printf("notes: getmarkdates: %v", err)
		return
	}
	fmt.Fprintf(w, "%s", string(jbody[:len(jbody)]))
}

func GetNotebooks(w http.ResponseWriter, r *http.Request) {
	_, err := app.GetSession(r)
	if err != nil {
		app.Exit(w, r)
		return
	}

	notebooks := AllNotes.GetNotebooks()
	// Write template
	tmpl := template.Must(template.ParseFiles(app.AppDir + "/model/notes/tmpl/notebooks.html"))
	if err := tmpl.Execute(w, notebooks); err != nil {
		log.Printf("%v", err)
		return
	}
}
