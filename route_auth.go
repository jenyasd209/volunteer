package main

import (
	"fmt"
	"graduate/data"
	"graduate/data/user"
	"graduate/data/user/freelancer"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("_cookie"); err == nil {
		http.Redirect(w, r, "/my_profile/about", 302)
	} else {
		data := &Data{"Login", nil}
		generateHTML(w, data, "base", "header", "footer", "login")
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("_cookie"); err == nil {
		http.Redirect(w, r, "/my_profile/about", 302)
	} else {
		data := &Data{"Registation", nil}
		generateHTML(w, data, "base", "header", "footer", "registration")
	}
}

func registrationAccount(w http.ResponseWriter, r *http.Request) {
	var user user.Userable
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	user = &freelancer.Freelancer{
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Password:  r.PostFormValue("password"),
		Email:     r.PostFormValue("email"),
	}

	if err := user.Create(); err != nil {
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
		session, err := freelancer.CreateSession()
		if err != nil {
			fmt.Println(err)
			return
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.GetUUID(),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/my_profile/about", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
