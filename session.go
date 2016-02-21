package main

import (
	"crypto/rand"
	"fmt"
	"time"
)

var sessionTable map[string]Session

func init() {
	sessionTable = make(map[string]Session)
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

func GetSession(key string) (*Session, err) {
	ok, s := sessionTable[key]
	if !ok {
		return nil, fmt.Errorf("No session for this key")
	}
	return s, nil
}
