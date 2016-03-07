package notes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"worgenda/app"

	"github.com/stacktic/dropbox"
)

const (
	TIMESYNCSOURCES = 1 * time.Minute
)

type DropboxConfig struct {
	AppKey    string
	AppSecret string
	Token     string
	Files     []string
}

func Sync() {
	dc, err := GetDropboxConfig()
	if err != nil {
		log.Printf("notes: sync: %v", err)
	}

	for {
		readSources(dc)
		time.Sleep(TIMESYNCSOURCES)
	}
}

func Upload() {
	dc, err := GetDropboxConfig()
	if err != nil {
		log.Printf("notes: upload: %v", err)
	}

	writeSources(dc)
}

func readSources(config *DropboxConfig) {
	for _, file := range config.Files {
		fcontent, err := ReadFile(config, file)
		if err != nil {
			log.Printf("notes: readsources: %v", err)
			continue
		}

		AllNotes.AddNotebook(filepath.Base(file), fcontent)
	}
}

func writeSources(config *DropboxConfig) {
	notesToWrite := AllNotes.GetNotesFromNotebook("worgenda.org")
	content := "#+TITLE: Worgenda Notebook\n\n"
	for _, note := range notesToWrite {
		if note.IsValid() {
			content += note.String()
		}
	}

	file := ""
	for _, file = range config.Files {
		if filepath.Base(file) == "worgenda.org" {
			break
		}
	}
	if file == "" {
		log.Printf("notes: writesource: no worgenda notebook found")
	}

	err := WriteFile(config, file, content)
	if err != nil {
		log.Printf("notes: writesources: %v", err)
	}
}

func GetDropboxConfig() (*DropboxConfig, error) {
	config := new(DropboxConfig)

	configFile, err := os.Open(app.AppDir + "/var/config.json")
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ReadFile(config *DropboxConfig, file string) (string, error) {

	var err error
	var db *dropbox.Dropbox
	var s string

	db = dropbox.NewDropbox()
	db.SetAppInfo(config.AppKey, config.AppSecret)
	db.SetAccessToken(config.Token)

	rd, _, err := db.Download(file, "", 0)
	if err != nil {
		return s, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rd)
	s = buf.String()

	return s, nil
}

func WriteFile(config *DropboxConfig, file string, content string) error {
	var db *dropbox.Dropbox

	db = dropbox.NewDropbox()
	db.SetAppInfo(config.AppKey, config.AppSecret)
	db.SetAccessToken(config.Token)

	in := ioutil.NopCloser(bytes.NewBufferString(content))
	_, err := db.UploadByChunk(in, 1024, file, true, "")
	return err
}
