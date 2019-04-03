package main

import (
	"graduate/data"
	"graduate/data/user/freelancer"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "base", "header", "footer", "login")
}

func registration(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "base", "header", "footer", "registration")
}

func registrationAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	freelancer := freelancer.Freelancer{
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Email:     r.PostFormValue("email"),
		Password:  r.PostFormValue("password"),
	}

	if err := freelancer.Create(); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/my_profile/about", 302)
}

func loginAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	freelancer, err := freelancer.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}

	if freelancer.Password == data.Encrypt(r.PostFormValue("password")) {
		// http.Redirect(w, r, "/my_profile/about", 302)
		generateHTML(w, &freelancer, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
