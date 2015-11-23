package main

import (
	//"io/ioutil"
	"regexp"
	"strings"
	"time"
)

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

// Recover cached notes (from datastore by now)
func LoadCachedNotes() ([]Note, error) {
	notes := make([]Note, 0)

	return notes, nil
}

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

	/*buf, err := ioutil.ReadFile(file)
	if err != nil {
		return notes, err
	}

	content := string(buf)*/

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
