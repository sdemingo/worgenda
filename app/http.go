package app

import (
	"log"
	"net/http"
	"text/template"
)

var AppDir string

func Run(appdir string) {
	AppDir = appdir

	LoadUsers()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
}

func Exit(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, r, http.StatusNotFound)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := GetSession(r)
	if err != nil {
		Exit(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("model/agenda/tmpl/agenda.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("%v", err)
		return
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		tmpl := template.Must(template.ParseFiles(AppDir + "/app/tmpl/error404.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("%v", err)
			return
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := GetSession(r)
	if s != nil {
		welcome(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles(AppDir + "/app/tmpl/login.html"))
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
	if err != nil || !user.PasswordOk(password) {
		log.Printf("Failed login for user %s", username)
		Exit(w, r)
		return
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
	Exit(w, r)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/notes/main", http.StatusFound)
}
