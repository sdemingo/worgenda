package notes

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"os"

	"github.com/stacktic/dropbox"
)

type DropboxConfig struct {
	AppKey    string
	AppSecret string
	Token     string
	Files     []string
}

func GetDropboxConfig() (*DropboxConfig, error) {
	config := new(DropboxConfig)

	configFile, err := os.Open("config.json")
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
