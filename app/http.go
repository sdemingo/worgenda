package app

import (
	"log"
	"net/http"
	"text/template"
)

func Run() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := GetSession(r)
	if s != nil {
		welcome(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("./app/tmpl/login.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := GetSession(r)
	if s != nil {
		welcome(w, r)
		return
	}

	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	user, err := GetUser(username)
	if err != nil || user.Password != password {
		exit(w, r)
	}

	log.Printf("User %s do login", user.Username)
	session := NewSession(user)
	sessionCookie := &http.Cookie{Name: "sessionKey", Value: session.Key, HttpOnly: false}
	http.SetCookie(w, sessionCookie)

	welcome(w, r)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := GetSession(r)
	if s != nil {
		log.Printf("User %s do login", s.User.Username)
	}
	DeleteSession(r)
	exit(w, r)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := GetSession(r)
	if err != nil {
		exit(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("model/agenda/tmpl/agenda.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}

func exit(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/welcome", http.StatusFound)
}
