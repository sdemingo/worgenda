package app

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
)

var userTable []*User

type User struct {
	Fullname string
	Username string
	Password string
}

func LoadUsers() error {
	var err error
	userTable = make([]*User, 0)

	usersDb, err := os.Open(AppDir + "/var/users.json")
	if err != nil {
		return err
	}

	usersJs := json.NewDecoder(usersDb)
	if err = usersJs.Decode(&userTable); err != nil {
		return err
	}
	return nil
}

func GetUser(key string) (*User, error) {
	for _, u := range userTable {
		if u.Username == key {
			return u, nil
		}
	}
	return nil, fmt.Errorf("app: user not found")
}

func (u *User) PasswordOk(pass string) bool {
	h := md5.New()
	h.Write([]byte(pass))
	hpass := fmt.Sprintf("%x", h.Sum(nil))
	return u.Password == hpass
}
