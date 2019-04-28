package main

import (
	"fmt"
	"graduate/data"
	"html/template"
	"net/http"
)

// var session data.SessionHelper = &data.Session{}
var session data.Session

func freelancerProfile(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/user/about", 302)
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
		generateHTML(w, data, funcMap, "base", "header", "footer", "userProfile/profile", "userProfile/about_base", "userProfile/freelancer/about")
	}
}

func freelancerProfileWorks(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		data := &Data{"My works", &freelancer}
		generateHTML(w, data, nil, "base", "header", "footer", "userProfile/profile", "userProfile/freelancer/my_works")
	}
}

func freelancerProfileContacts(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		data := &Data{"Contacts", &freelancer}
		generateHTML(w, data, nil, "base", "header", "footer", "userProfile/profile", "userProfile/contacts")
	}
}

func freelancerProfileSetting(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		specs, _ := data.GetAllSpecialization()
		if r.Method == http.MethodPost {
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
			freelancer.Specialization = arrayStringToArrayInt(r.Form["specialization[]"])
			err = freelancer.Update()
			if err != nil {
				fmt.Println(err)
			}
			http.Redirect(w, r, "/user/about", 302)
		} else {
			type Content struct {
				User           data.Freelancer
				Specialization []data.Specialization
			}
			funcMap := template.FuncMap{
				"getNameSpecialization":  data.GetSpecializationName,
				"containsSpecialization": freelancer.ContainsSpecialization,
			}
			data := &Data{"Setting", &Content{freelancer, specs}}
			generateHTML(w, data, funcMap, "base", "header", "footer", "userProfile/profile", "userProfile/setting_base", "userProfile/freelancer/setting")
		}
	}
}
