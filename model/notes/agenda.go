package notes

import (
	"sync"
	"time"
)

var AllNotes *Agenda

type Agenda struct {
	Notebooks map[string]string
	Notes     map[string][]Note
	rMutex    sync.RWMutex
	lastSync  time.Time
}

func init() {
	AllNotes = NewAgenda()
}

func NewAgenda() *Agenda {
	a := new(Agenda)
	a.Notes = make(map[string][]Note)
	a.Notebooks = make(map[string]string)
	return a
}

func (a *Agenda) GetNote(id int64) *Note {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	for _, notes := range a.Notes {
		for i := range notes {
			if notes[i].Id == id {
				return &notes[i]
			}
		}
	}
	return nil
}

/*
func (a *Agenda) Build(notes []Note) error {
	a.rMutex.RLock()
	defer a.rMutex.RUnlock()

	a.Notes = notes
	a.lastSync = time.Now()
	return nil
}
*/
func (a *Agenda) AddNotebook(name string, content string) error {
	a.rMutex.RLock()
	defer a.rMutex.RUnlock()

	notes := Parse(content)
	a.Notes[name] = make([]Note, 0, len(notes))
	for i := range notes {
		notes[i].Source = name
		a.Notes[name] = append(a.Notes[name], notes[i])
	}

	a.Notebooks[name] = content
	a.lastSync = time.Now()

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

func (a *Agenda) GetNotesFromDate(daynotes *DayNotes) {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	for _, notes := range a.Notes {
		for i := range notes {
			daynotes.Add(notes[i])
		}
	}
}

func (a *Agenda) GetBusyDates() []string {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	dates := make(map[time.Time]bool)
	for _, notes := range a.Notes {
		for _, note := range notes {
			for _, stamp := range note.Stamps {
				if ok, _ := dates[stamp]; !ok {
					dates[stamp] = true
				}
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
