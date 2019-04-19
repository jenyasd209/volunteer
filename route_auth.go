package main

import (
	"fmt"
	"graduate/data"
	"log"
	"net/http"
)

var pageData *Data

func login(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("_cookie"); err == nil {
		http.Redirect(w, r, "/my_profile/about", 302)
	} else {
		pageData = &Data{"Login", nil}
		generateHTML(w, pageData, "base", "header", "footer", "login")
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("_cookie"); err == nil {
		http.Redirect(w, r, "/my_profile/about", 302)
	} else {
		pageData := &Data{"Registation", nil}
		generateHTML(w, pageData, "base", "header", "footer", "registration")
	}
}

func registrationAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	group := r.Form.Get("group")
	user := &data.User{
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		Password:  r.PostFormValue("password"),
		Email:     r.PostFormValue("email"),
		Phone:     r.PostFormValue("phone"),
		Skype:     r.PostFormValue("skype"),
		Facebook:  r.PostFormValue("facebook"),
	}
	if group == "volunteer" {
		specialization := arrayStringToArrayInt(r.Form["specialization[]"])
		freelancer := &data.Freelancer{
			Specialization: specialization,
			User:           *user,
		}
		if err := freelancer.Create(); err != nil {
			fmt.Println(err)
		}
	}
	if group == "customer" {
		customer := &data.Customer{
			Organization: r.PostFormValue("organization-name"),
			User:         *user,
		}
		if err := customer.Create(); err != nil {
			fmt.Println(err)
		}
	}

	http.Redirect(w, r, "/my_profile/about", 302)
}

func loginAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	user, err := data.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", 404)
	}

	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		group := r.Form.Get("group")
		if group == "volunteer" {
			if ok, _ := data.CheckFreelancer(user.ID); !ok {
				http.Redirect(w, r, "/login", 302)
				return
			}
		} else if group == "customer" {
			if ok, _ := data.CheckCustomer(user.ID); !ok {
				http.Redirect(w, r, "/login", 302)
				return
			}
		} else {
			http.Redirect(w, r, "/login", 302)
			return
		}
		session, err := user.CreateSession()
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

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		log.Println(err, " Failed to get cookie")
		session := data.Session{UUID: cookie.Value}
		log.Println(session)
		session.DeleteByUUID()
		c := http.Cookie{
			Name:   "_cookie",
			MaxAge: -1}
		http.SetCookie(w, &c)
	}
	http.Redirect(w, r, "/", 302)
}
