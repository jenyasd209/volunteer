package main

import (
	"fmt"
	"graduate/data"
	"net/http"
)

// var session data.SessionHelper = &data.Session{}
var session data.Session

func freelancerProfile(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/my_profile/about", 302)
}

func freelancerProfileAbout(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// freelancer, err := data.GetFreelancerByUserID(session.GetUserID())
		freelancer, err := data.GetFreelancerByUserID(session.UserID)
		if err != nil {
			fmt.Println(err)
		}
		data := &Data{"About", &freelancer}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
	}
}

func freelancerProfileWorks(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// freelancer, _ := data.GetFreelancerByUserID(session.GetUserID())
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		data := &Data{"My works", &freelancer}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/my_works")
	}
}

func freelancerProfileContacts(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// freelancer, _ := data.GetFreelancerByUserID(session.GetUserID())
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		data := &Data{"Contacts", &freelancer}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/contacts")
	}
}
