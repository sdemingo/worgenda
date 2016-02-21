package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

var sessionTable map[string]*Session

func init() {
	sessionTable = make(map[string]*Session)
}

type Session struct {
	Key       string
	User      *User
	LoginTime time.Time
}

func randKey() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func NewSession(u *User) *Session {
	key := randKey()
	s := &Session{key, u, time.Now()}
	sessionTable[key] = s
	return s
}

func GetSession(r *http.Request) (*Session, error) {

	ck, err := r.Cookie("sessionKey")
	if err != nil {
		return nil, err
	}

	s, ok := sessionTable[ck.Value]
	if !ok {
		return nil, fmt.Errorf("No session for this key")
	}
	return s, nil
}

func DeleteSession(r *http.Request) error {

	s, err := GetSession(r)
	if err != nil {
		return err
	}
	delete(sessionTable, s.Key)
	return nil
}
