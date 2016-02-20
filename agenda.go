package main

import (
	//"fmt"
	"time"
)

const (
	NEXT_DAYS           = 20
	DATEHOURFORMATPRINT = "02/01/2006 Mon 15:04"
	DATEFORMATPRINT     = "02/01/2006 Mon"
)

type Day struct {
	Tm       time.Time
	TmStr    string
	SchedDay []Note
}

type Agenda struct {
	AllNotes []Note
	Sched    []Day
}

func NewAgenda() *Agenda {
	a := new(Agenda)
	return a
}

func (a *Agenda) InsertNotes(notes []Note) {
	for i := range notes {
		a.AllNotes = append(a.AllNotes, notes[i])
	}
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
}

func (a *Agenda) String() string {
	s := ""
	for _, day := range a.Sched {
		s = s + " [" + day.Tm.Format(DATEHOURFORMATPRINT) + "]\n"
		for _, note := range day.SchedDay {
			s = s + "\t -" + note.Title + "\n"
		}
		s = s + "\n"
	}
	return s
}
