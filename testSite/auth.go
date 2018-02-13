package main

import (
	"net/http"
	"time"
)

const EXP_TIMER = 120

func authIssue(w http.ResponseWriter, uVal string) {
	exp := time.Now().Add(time.Minute * EXP_TIMER)
	cookie := http.Cookie{Name: "username", Value: uVal, Expires: exp}
	http.SetCookie(w, &cookie)
}

func authCheck(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		authIssue(w, cookie.Value)
	}
}
