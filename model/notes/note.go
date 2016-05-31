package notes

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const (
	DATEFORMATGETPARAM  = "02/01/2006"
	DATEHOURFORMATPRINT = "02/01/2006 Mon 15:04"
	HOURFORMATPRINT     = "15:04"
	DATEFORMATPRINT     = "02 Jan 2006"
	DATEFORMATFORHTML   = "Monday, 02 January 2006"

	MAXWORDSRESUMEBODY = 20
)

var nullTime = time.Time{}

type Note struct {
	Id       int64
	Title    string
	Body     string
	Stamps   []time.Time
	Source   string
	Status   string
	Deadline time.Time
	Warning  time.Duration
}

func NewNote() *Note {
	n := new(Note)
	n.Id = rand.Int63()
	n.Stamps = make([]time.Time, 0)
	n.Status = ""
	n.Deadline = nullTime
	n.Warning = time.Duration(0)
	return n
}

func (n *Note) IsValid() bool {
	return len(strings.TrimSpace(n.Title)) > 0
}

func (n *Note) IsTodo() bool {
	return n.Status == "TODO"
}

func (n *Note) HasDeadline(date time.Time) bool {
	if n.Deadline == nullTime {
		return false
	}
	return n.Deadline == date
}

func (n *Note) IsWarningTime() bool {
	if n.Warning == 0 || n.Deadline == nullTime {
		return true
	}
	warnTime := time.Now().Add(n.Warning)
	return warnTime.After(n.Deadline)
}

func (n *Note) GetResumeBody() string {
	words := strings.Fields(n.Body)
	if len(words) < MAXWORDSRESUMEBODY {
		return strings.Join(words, " ")
	} else {
		return ""
	}
}

func (n *Note) GetTextBody() string {
	body := Org2HTML([]byte(n.Body), "")
	return body
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

// Return the string part with the date for the timestamp that happens
// in this day
func (n *Note) GetStampDate(date time.Time) string {
	stamp := n.GetStampForDay(date)
	return stamp.Format(DATEFORMATFORHTML)
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

func (n *Note) String() string {
	s := fmt.Sprintf("\n* %s %s\n", n.Status, n.Title)
	for i := range n.Stamps {
		s += fmt.Sprintf("\t<%s>\n", n.Stamps[i].Format(ORGDATEHOURFORMAT))
	}
	s += n.Body + "\n"
	return s
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

type Bookmark struct {
	Url  string
	Text string
	Desc string
	Tags []string
}

func NewBookmark(note Note) Bookmark {
	b := new(Bookmark)
	b.Tags = make([]string, 0)

	re := regexp.MustCompile("\\[\\[.+\\]\\]")
	links := re.FindAllString(note.Title, 1)
	if len(links) > 0 {
		fields := strings.Split(strings.Trim(links[0], "[]"), "][")
		b.Url = fields[0]
		b.Text = fields[1]
	}

	tagsRe := regexp.MustCompile(":[a-zA-Z0-9:]+:")
	tagsFind := tagsRe.FindAllString(note.Title, -1)
	if len(tagsFind) > 0 {
		b.Tags = strings.Split(tagsFind[0], ":")
	}
	return *b
}
