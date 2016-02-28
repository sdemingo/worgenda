package notes

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
	//"fmt"
	"os"
	"path/filepath"

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
		log.Panic(err)
	}

	for {
		readSources(dc)
		time.Sleep(TIMESYNCSOURCES)
	}
}

func readSources(config *DropboxConfig) {
	newNotes := make([]Note, 0)
	for _, file := range config.Files {
		fcontent, err := ReadFile(config, file)
		if err != nil {
			log.Printf("notes: sync: %v", err)
			continue
		}

		notes := Parse(fcontent)
		for i := range notes {
			notes[i].Source = filepath.Base(file)
			newNotes = append(newNotes, notes[i])
		}
	}

	AllNotes = newNotes
	log.Println("Sources was syncing succesfully")
}

func GetDropboxConfig() (*DropboxConfig, error) {
	config := new(DropboxConfig)

	configFile, err := os.Open("var/config.json")
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
