package notes

import (
	"log"
	"time"
)

const (
	DATEHOURFORMATPRINT = "02/01/2006 Mon 15:04"
	HOURFORMATPRINT     = "15:04"
	DATEFORMATPRINT     = "02 Jan 2006"
	DATEFORMATFORHTML   = "Monday, 02 January 2006"
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

// Return the string part with the hour for the timestamp that happens
// in this day
func (n *Note) GetStampHour(date time.Time) string {
	stamp := n.GetStampForDay(date)
	if stamp.Hour() == 0 && stamp.Minute() == 0 {
		return ""
	}
	return stamp.Format(HOURFORMATPRINT)
}

// Return if a note have a stamp which happens in this day
func (n *Note) InDay(date time.Time) bool {
	for _, stamp := range n.Stamps {
		md := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 1, 0, time.UTC)
		diffMd := int(stamp.Sub(md).Hours())
		if stamp == date || (diffMd < 24 && diffMd > 0) {
			return true
		}
	}
	return false
}

// Get the stamp of the note for this day
func (n *Note) GetStampForDay(date time.Time) time.Time {
	for _, stamp := range n.Stamps {
		md := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 1, 0, time.UTC)
		diffMd := int(stamp.Sub(md).Hours())
		if stamp == date || (diffMd < 24 && diffMd > 0) {
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
