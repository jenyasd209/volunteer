package main

import (
	"fmt"
	"graduate/data"
	siteUser "graduate/data/user"
	"graduate/data/user/freelancer"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	data := &Data{"Login", nil}
	generateHTML(w, data, "base", "header", "footer", "login")
}

func registration(w http.ResponseWriter, r *http.Request) {
	data := &Data{"Registation", nil}
	generateHTML(w, data, "base", "header", "footer", "registration")
}

func registrationAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	freelancer := freelancer.Freelancer{
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Password:  r.PostFormValue("password"),
		UserType: siteUser.UserType{
			Email: r.PostFormValue("email"),
		},
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
		fmt.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	if freelancer.Password == data.Encrypt(r.PostFormValue("password")) {
		siteUser.User.Set(freelancer.Email, freelancer.ID, true)
		freelancer.CreateSession()
		generateHTML(w, nil, "base", "header", "footer", "userProfile/worker_personal_profile", "userProfile/about")
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
