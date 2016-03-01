package notes

import (
	"sync"
	"time"
)

var AllNotes *Agenda

type Agenda struct {
	Notebooks map[string]string
	Notes     []Note
	rMutex    sync.RWMutex
	lastSync  time.Time
}

func init() {
	AllNotes = NewAgenda()
}

func NewAgenda() *Agenda {
	a := new(Agenda)
	a.Notes = make([]Note, 0)
	a.Notebooks = make(map[string]string)
	return a
}

func (a *Agenda) GetNote(id int64) *Note {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	for i := range a.Notes {
		if a.Notes[i].Id == id {
			return &a.Notes[i]
		}
	}
	return nil
}

func (a *Agenda) Build(notes []Note) error {
	a.rMutex.RLock()
	defer a.rMutex.RUnlock()

	a.Notes = notes
	a.lastSync = time.Now()
	return nil
}

func (a *Agenda) AddNotebook(name string, content string) error {
	a.rMutex.RLock()
	defer a.rMutex.RUnlock()

	notes := Parse(content)
	for i := range notes {
		notes[i].Source = name
		a.Notes = append(a.Notes, notes[i])
	}

	a.Notebooks[name] = content

	return nil
}

func (a *Agenda) GetNotebooks() map[string]string {
	notebooks := make(map[string]string)
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	for k, v := range a.Notebooks {
		notebooks[k] = v
	}
	return notebooks
}

func (a *Agenda) GetNotesFromDate(notes *DayNotes) {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()
	for _, note := range a.Notes {
		notes.Add(note)
	}
}

func (a *Agenda) GetBusyDates() []string {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	dates := make(map[time.Time]bool)
	for _, note := range a.Notes {
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

	return alldates
}
