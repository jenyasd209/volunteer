package main

import (
	"fmt"
	"graduate/data/user/freelancer"
	"net/http"
)

func freelancerProfile(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/my_profile/about", 302)
}
func freelancerProfileAbout(w http.ResponseWriter, r *http.Request) {
	s := freelancer.Session{}
	sess, err := session(w, r, &s)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, err := sess.GetUser()
		fmt.Println(sess)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
		data := &Data{"About", &user}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
	}
}

func freelancerProfileWorks(w http.ResponseWriter, r *http.Request) {
	s := freelancer.Session{}
	sess, err := session(w, r, &s)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, _ := sess.GetUser()
		data := &Data{"My works", &user}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/my_works")
	}
}

func freelancerProfileContacts(w http.ResponseWriter, r *http.Request) {
	s := freelancer.Session{}
	sess, err := session(w, r, &s)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, _ := sess.GetUser()
		data := &Data{"Contacts", &user}
		generateHTML(w, data, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/contacts")
	}
}
