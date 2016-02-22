package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var userTable []*User

func init() {
	var err error
	userTable = make([]*User, 0)

	usersDb, err := os.Open("./var/users.json")
	if err != nil {
		log.Panic(err)
	}

	usersJs := json.NewDecoder(usersDb)
	if err = usersJs.Decode(&userTable); err != nil {
		log.Panic(err)
	}
}

type User struct {
	Fullname string
	Username string
	Password string
}

func GetUser(key string) (*User, error) {
	for _, u := range userTable {
		if u.Username == key {
			return u, nil
		}
	}
	return nil, fmt.Errorf("app: user not found")
}
