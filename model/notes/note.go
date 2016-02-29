package notes

import (
	"math/rand"
	"strings"
	"time"
)

const (
	DATEHOURFORMATPRINT = "02/01/2006 Mon 15:04"
	HOURFORMATPRINT     = "15:04"
	DATEFORMATPRINT     = "02 Jan 2006"
	DATEFORMATFORHTML   = "Monday, 02 January 2006"

	MAXWORDSRESUMEBODY = 20
)

type Note struct {
	Id     int64
	Title  string
	Body   string
	Stamps []time.Time
	Source string
}

func NewNote() *Note {
	n := new(Note)
	n.Id = rand.Int63()
	return n
}

func (n *Note) GetResumeBody() string {
	words := strings.Fields(n.Body)
	if len(words) < MAXWORDSRESUMEBODY {
		return strings.Join(words, " ")
	} else {
		return strings.Join(words[:MAXWORDSRESUMEBODY], " ") + " ..."
	}
}

func (n *Note) GetHTMLBody() string {
	return Org2HTML([]byte(n.Body), "")
}

// Return if a note have a stamp which happens in this day
func (n *Note) InDay(date time.Time) bool {
	for _, stamp := range n.Stamps {
		md := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 1, 0, time.UTC)
		diffMd := int(stamp.Sub(md).Hours())
		if stamp == date || (diffMd < 24 && diffMd >= 0) {
			return true
		}
	}
	return false
}

// Return the string part with the hour for the timestamp that happens
// in this day
func (n *Note) GetStampHour(date time.Time) string {
	stamp := n.GetStampForDay(date)
	if stamp.Hour() == 0 && stamp.Minute() == 0 {
		return ""
	}
	return stamp.Format(HOURFORMATPRINT)
}

// Get the stamp of the note for this day
func (n *Note) GetStampForDay(date time.Time) time.Time {
	for _, stamp := range n.Stamps {
		if datesInSameDay(stamp, date) {
			return stamp
		}
	}
	return time.Date(1982, time.August, 20, 20, 45, 0, 0, time.UTC)
}

// DayNotes sort a group of notes for a particular day
type DayNotes struct {
	Date  time.Time
	Notes []Note
}

func NewDayNotes(date time.Time) *DayNotes {
	dn := new(DayNotes)
	dn.Notes = make([]Note, 0)
	dn.Date = date
	return dn
}

func (dn *DayNotes) Add(note Note) {
	if note.InDay(dn.Date) {
		dn.Notes = append(dn.Notes, note)
	}
}

func (dn *DayNotes) Len() int {
	return len(dn.Notes)
}

func (dn *DayNotes) Less(i, j int) bool {
	return dn.Notes[i].GetStampForDay(dn.Date).Before(dn.Notes[j].GetStampForDay(dn.Date))
}

func (dn *DayNotes) Swap(i, j int) {
	dn.Notes[i], dn.Notes[j] = dn.Notes[j], dn.Notes[i]
}

func datesInSameDay(date1, date2 time.Time) bool {
	var md time.Time
	var diffMd int
	if date1.Before(date2) {
		md = time.Date(date1.Year(), date1.Month(), date1.Day(), 0, 0, 1, 0, time.UTC)
		diffMd = int(date2.Sub(md).Hours())
	} else {
		md = time.Date(date2.Year(), date2.Month(), date2.Day(), 0, 0, 1, 0, time.UTC)
		diffMd = int(date1.Sub(md).Hours())
	}
	if date1 == date2 || (diffMd < 24 && diffMd >= 0) {
		return true
	}

	return false
}
