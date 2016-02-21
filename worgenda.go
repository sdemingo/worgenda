package main

import (
	"html/template"
	"log"
	"net/http"
)

const (
	TLSPORT    = ":8443"
	PORT       = ":8080"
	PRIV_KEY   = "./keys/private_key"
	PUBLIC_KEY = "./keys/public_key"
	DOMAIN     = "localhost"
)

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/welcome", welcomeHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Redirect all requests to TLS socket
	go func() {
		err := http.ListenAndServe(PORT, http.RedirectHandler("https://"+DOMAIN+TLSPORT, http.StatusFound))
		if err != nil {
			panic("Error: " + err.Error())
		}
	}()

	// Listen on secure port
	err := http.ListenAndServeTLS(TLSPORT, PUBLIC_KEY, PRIV_KEY, nil)
	if err != nil {
		panic("Error: " + err.Error())
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("login.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	user, err := GetUser(username)
	if err != nil || user.Password != password {
		exit(w, r)
	}

	welcome(w, r)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("model/agenda/tmpl/agenda.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}

func exit(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/welcome", http.StatusMovedPermanently)
}
