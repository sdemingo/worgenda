package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"
)

var sessionTable map[string]*Session

const (
	TIME_MAX_INACTIVE_SESSION  = 2 * time.Minute
	TIME_SESSIONCLEANER_PERIOD = 1 * time.Minute
)

func init() {
	sessionTable = make(map[string]*Session)
	go sessionCleaner()
}

type Session struct {
	Key       string
	User      *User
	LoginTime time.Time
	LastTime  time.Time
}

func sessionCleaner() {
	log.Printf("Run session cleaner")
	for {
		for _, s := range sessionTable {
			if time.Since(s.LastTime) > TIME_MAX_INACTIVE_SESSION {
				delete(sessionTable, s.Key)
				log.Printf("Session for user %s closed for inactivity", s.User.Username)
			}
		}
		time.Sleep(TIME_SESSIONCLEANER_PERIOD)
	}
}

func randKey() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func NewSession(u *User) *Session {
	key := randKey()
	s := &Session{key, u, time.Now(), time.Now()}
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

	s.LastTime = time.Now()
	return s, nil
}

func DeleteSession(r *http.Request) error {

	s, err := GetSession(r)
	if err != nil {
		return err
	}
	delete(sessionTable, s.Key)
	log.Printf("Session for user %s closed", s.User.Username)
	return nil
}
