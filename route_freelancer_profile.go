package main

import (
	"fmt"
	"graduate/data"
	"html/template"
	"log"
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
		freelancer, err := data.GetFreelancerByUserID(session.UserID)
		if err != nil {
			fmt.Println(err)
		}
		funcMap := template.FuncMap{
			"getNameSpecialization": data.GetSpecializationName,
		}
		data := &Data{"About", &freelancer}
		generateHTML(w, data, funcMap, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
	}
}

func freelancerProfileWorks(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		data := &Data{"My works", &freelancer}
		generateHTML(w, data, nil, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/my_works")
	}
}

func freelancerProfileContacts(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		data := &Data{"Contacts", &freelancer}
		generateHTML(w, data, nil, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/contacts")
	}
}

func freelancerProfileSetting(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		if r.Method == http.MethodPost {
			log.Println(r.Method)
			err := r.ParseForm()
			if err != nil {
				panic(err)
			}
			freelancer.User.FirstName = r.PostFormValue("first_name")
			freelancer.User.LastName = r.PostFormValue("last_name")
			freelancer.User.About = r.PostFormValue("about")
			freelancer.User.Phone = r.PostFormValue("phone")
			freelancer.User.Facebook = r.PostFormValue("facebook")
			freelancer.User.Skype = r.PostFormValue("skype")
			err = freelancer.Update()
			if err != nil {
				fmt.Println(err)
			}
			http.Redirect(w, r, "/my_profile/about", 302)
		} else {
			type Content struct {
				User data.Freelancer
			}
			data := &Data{"Setting", &Content{freelancer}}
			generateHTML(w, data, nil, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/setting")
		}
	}
}
