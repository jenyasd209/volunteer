package main

import (
	"fmt"
	"graduate/data/user"
	"net/http"
)

// var sess user.SessionHelper = &user.Session{}
var session user.SessionHelper = &user.Session{}

func freelancerProfile(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/my_profile/about", 302)
}
func freelancerProfileAbout(w http.ResponseWriter, r *http.Request) {
	err := user.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, err := user.GetByUserID(session.GetUserID())
		if err != nil {
			fmt.Println(err)
		}
		data := &Data{"About", &freelancer}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
	}
}

func freelancerProfileWorks(w http.ResponseWriter, r *http.Request) {
	err := user.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := user.GetByUserID(session.GetUserID())
		data := &Data{"My works", &freelancer}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/my_works")
	}
}

func freelancerProfileContacts(w http.ResponseWriter, r *http.Request) {
	err := user.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := user.GetByUserID(session.GetUserID())
		data := &Data{"Contacts", &freelancer}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/contacts")
	}
}
