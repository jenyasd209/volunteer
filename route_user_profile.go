package main

import (
	"net/http"
)

func freelancerProfile(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/my_profile/about", 302)
}

func freelancerProfileAbout(w http.ResponseWriter, r *http.Request) {
	data := &Data{"About", nil}
	generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
}

func freelancerProfileWorks(w http.ResponseWriter, r *http.Request) {
	data := &Data{"My works", nil}
	generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/my_works")
}

func freelancerProfileContacts(w http.ResponseWriter, r *http.Request) {
	data := &Data{"Contacts", nil}
	generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/contacts")
}
