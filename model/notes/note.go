package notes

import (
	//"io/ioutil"

	"log"
	"regexp"
	"strings"
	"time"
)

const (
	DATEHOURFORMATPRINT = "02/01/2006 Mon 15:04"
	DATEFORMATPRINT     = "02 Jan 2006"
)

var AllNotes []Note

func init() {
	// Load notes from source (dropbox)

	dc, err := GetDropboxConfig()
	if err != nil {
		log.Panic(err)
	}
	content, err := ReadFile(dc, "prueba.org")
	if err != nil {
		log.Panic(err)
	}
	AllNotes = Parse(content)
	log.Printf("All notes are loaded")
}

type Note struct {
	Title  string
	Body   string
	Stamps []time.Time
}

func (n *Note) String() string {
	s := n.Title + "\n"
	s = s + n.Body + "\n"
	for _, st := range n.Stamps {
		s = s + "\t -" + st.Format(DATEHOURFORMATPRINT) + "\n"
	}
	s = s + "\n\n"
	return s
}

/*
// Recover cached notes (from datastore by now)
func LoadCachedNotes() ([]Note, error) {
	notes := make([]Note, 0)

	return notes, nil
}
*/

/*

   Org mode

*/

var noteTitleReg = regexp.MustCompile("(?m)^(\\*{1,3} .+\\n)")
var separator = "@@@@\n"
var separatorReg = regexp.MustCompile("(?m)^@@@@\\n")
var dateReg = regexp.MustCompile("\\<\\d{4}-\\d{2}-\\d{2} .{3}( \\d{2}:\\d{2})?\\>")

// Parse string content in org mode and recover notes from it
func Parse(content string) []Note {
	notes := make([]Note, 0)

	content = noteTitleReg.ReplaceAllString(content, separator+"$1")
	rawNotes := separatorReg.Split(content, -1)

	for _, rnote := range rawNotes {
		note := new(Note)
		note.Title = parseTitle(rnote)
		note.Body = parseBody(rnote)
		note.Stamps = parseStamps(rnote)

		notes = append(notes, *note)
	}

	return notes
}

func parseTitle(orgnote string) string {
	title := noteTitleReg.FindString(orgnote)
	prefix := regexp.MustCompile("(?m)^\\*+")
	return strings.Trim(prefix.ReplaceAllString(title, ""), " \n\t")
}

func parseBody(orgnote string) string {
	body := noteTitleReg.ReplaceAllString(orgnote, "")
	body = dateReg.ReplaceAllString(body, "")
	return strings.Trim(body, " \n\t")
}

func parseStamps(orgnote string) []time.Time {

	times := make([]time.Time, 0)
	rawTimes := dateReg.FindAllString(orgnote, -1)

	for _, rt := range rawTimes {
		r := regexp.MustCompile(" [a-záéíóú]{3}")
		rt = r.ReplaceAllString(rt, "")
		var t time.Time
		var err error
		if strings.Contains(rt, ":") {
			t, err = time.Parse("<2006-01-02 15:04>", rt)
		} else {
			t, err = time.Parse("<2006-01-02>", rt)
		}
		if err == nil {
			times = append(times, t)
		}
	}

	return times
}

/*

Old agenda Code
===============


func NewAgenda() *Agenda {
	a := new(Agenda)
	return a
}


func (a *Agenda) FilterNotes(date time.Time) []Note {
	filter := make([]Note, 0)
	for _, note := range a.AllNotes {
		for _, stamp := range note.Stamps {
			day := date
			past := date.AddDate(0, 0, -1)

			if past.Before(stamp) && day.After(stamp) {
				filter = append(filter, note)
			}
		}
	}
	return filter
}

func (a *Agenda) Build() {
	now := time.Now()
	basetime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	c := 0
	a.Sched = make([]Day, NEXT_DAYS)
	for c < NEXT_DAYS {
		a.Sched[c].Tm = basetime.AddDate(0, 0, c)
		a.Sched[c].TmStr = a.Sched[c].Tm.Format(DATEFORMATPRINT)
		a.Sched[c].SchedDay = a.FilterNotes(a.Sched[c].Tm)
		c++
	}
}*/
