package main

import (
	"log"
	"net/http"
	"worgenda/app"
	"worgenda/model/notes"
)

const (
	TLSPORT    = ":8443"
	PORT       = ":8080"
	PRIV_KEY   = "./var/private_key"
	PUBLIC_KEY = "./var/public_key"
	DOMAIN     = "192.168.1.107"
)

func main() {
	log.Printf("Run worgenda")

	app.Run()

	http.HandleFunc("/welcome", notes.Main)
	http.HandleFunc("/notes/dates", notes.GetMarkDates)
	http.HandleFunc("/notes/events", notes.GetEvents)
	http.HandleFunc("/notes/event", notes.GetEvent)
	http.HandleFunc("/notes/new", notes.NewEventForm)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Redirect all requests to TLS socket
	go func() {
		err := http.ListenAndServe(PORT, http.RedirectHandler("https://"+DOMAIN+TLSPORT, http.StatusFound))
		if err != nil {
			panic("Error: " + err.Error())
		}
	}()

	// Run sync goroutine
	go notes.Sync()

	// Listen on secure port
	err := http.ListenAndServeTLS(TLSPORT, PUBLIC_KEY, PRIV_KEY, nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
