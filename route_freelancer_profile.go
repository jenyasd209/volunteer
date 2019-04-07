package main

import (
	"fmt"
	"graduate/data/freelancer"
	"graduate/user"
	"net/http"
)

var sess user.SessionHelper = &freelancer.Session{}

func freelancerProfile(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/my_profile/about", 302)
}
func freelancerProfileAbout(w http.ResponseWriter, r *http.Request) {
	err := user.Session(r, &sess)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			fmt.Println(err)
		}
		data := &Data{"About", &user}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
	}
}

func freelancerProfileWorks(w http.ResponseWriter, r *http.Request) {
	err := user.Session(r, &sess)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, _ := sess.User()
		data := &Data{"My works", &user}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/my_works")
	}
}

func freelancerProfileContacts(w http.ResponseWriter, r *http.Request) {
	err := user.Session(r, &sess)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, _ := sess.User()
		data := &Data{"Contacts", &user}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/contacts")
	}
}
