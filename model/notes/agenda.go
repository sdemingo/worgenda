package notes

import (
	"strings"
	"sync"
	"time"
)

var AllNotes *Agenda

type Agenda struct {
	Notebooks map[string]string
	Notes     map[string][]Note
	Bookmarks []Bookmark
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
	a.Bookmarks = make([]Bookmark, 0)
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

func (a *Agenda) AddNotebook(name string, content string) error {
	a.rMutex.RLock()
	defer a.rMutex.RUnlock()

	notes := Parse(content)
	a.Notes[name] = make([]Note, 0, len(notes))
	for i := range notes {
		notes[i].Source = name
		a.Notes[name] = append(a.Notes[name], notes[i])
	}

	a.Notebooks[name] = Org2HTML([]byte(content), "")
	a.lastSync = time.Now()

	return nil
}

func (a *Agenda) AddBookmarks(content string) error {
	a.rMutex.RLock()
	defer a.rMutex.RUnlock()

	notes := Parse(content)
	for i := range notes {
		b := NewBookmark(notes[i])
		if b.Text != "" {
			a.Bookmarks = append(a.Bookmarks, b)
		}

	}

	return nil
}

func (a *Agenda) AddNote(n Note) {
	a.rMutex.RLock()
	defer a.rMutex.RUnlock()

	if !n.IsValid() {
		return
	}

	if _, ok := a.Notes[n.Source]; !ok {
		a.Notes[n.Source] = make([]Note, 0)
	}

	a.Notes[n.Source] = append(a.Notes[n.Source], n)
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

func (a *Agenda) GetNotesWithText(pattern string) []Note {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	search := make([]Note, 0)
	pattern = strings.ToLower(pattern)

	for _, nbook := range a.Notes {
		for _, note := range nbook {
			title := strings.ToLower(note.Title)
			body := strings.ToLower(note.Body)
			if strings.Contains(title, pattern) ||
				strings.Contains(body, pattern) {
				search = append(search, note)
			}
		}
	}

	return search
}

func (a *Agenda) GetNotesFromNotebook(notebook string) []Note {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	if _, ok := a.Notes[notebook]; !ok {
		return nil
	}

	return a.Notes[notebook]
}

func (a *Agenda) GetNotesToDo() []Note {
	a.rMutex.Lock()
	defer a.rMutex.Unlock()

	todo := make([]Note, 0)

	for _, notes := range a.Notes {
		for _, note := range notes {
			if note.IsTodo() && note.IsWarningTime() {
				todo = append(todo, note)
			}
		}
	}

	return todo
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
