package main

import (
	"fmt"
	"graduate/data"
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	//if _, err := r.Cookie("_cookie"); err == nil {
	//	http.Redirect(w, r, "/my_profile", 302)
	_, err := session(w, r)
	if err == nil {
		http.Redirect(w, r, "/my_profile", 302)
	} else {
		pageData := PageData{
			Title :"Login",
		}
		//pageData.Title = "Login"
		generateHTML(w, &pageData, nil, "base", "header", "footer", "login")
	}
}

func registration(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("_cookie"); err == nil {
		http.Redirect(w, r, "/my_profile", 302)
	} else {
		specialization, _ := data.GetAllSpecialization()
		pageData := PageData{
			Title :"Registration",
			Content : struct{
				Specialization []data.Specialization
			}{specialization},
		}
		generateHTML(w, &pageData, nil, "base", "header", "footer", "registration")
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
		user.RoleID = data.UserRoleFreelancer
		freelancer := &data.Freelancer{
			Specialization: specialization,
			User:           *user,
		}
		if err := freelancer.Create(); err != nil {
			fmt.Println(err)
		}
	}
	if group == "customer" {
		user.RoleID = data.UserRoleCustomer
		customer := &data.Customer{
			Organization: r.PostFormValue("organization-name"),
			User:         *user,
		}
		if err := customer.Create(); err != nil {
			fmt.Println(err)
		}
	}

	http.Redirect(w, r, "/login", 302)
}

func loginAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	pageData := PageData{
		Title :"Registration",
	}

	user := data.User{
		Email: r.PostFormValue("email"),
		Password: data.Encrypt(r.PostFormValue("password")),
	}

	if user.CheckLoginData(){
		sess, err := user.CreateSession()
		if err != nil {
			return
		}
		cookie := http.Cookie{
			Name: "_cookie",
			Value:    sess.UUID,
			HttpOnly: true,
			MaxAge: 60 * 60 * 24 * 30,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/my_profile", 302)
	}else {
		pageData.Errors = []string{"User is not found"}
		generateHTML(w, &pageData, nil, "base", "header", "footer", "login")
		return
	}
	//user, err := data.GetUserByEmail(r.PostFormValue("email"))
	//if err != nil {
	//	pageData.Errors = []string{"User is not found"}
	//	generateHTML(w, &pageData, nil, "base", "header", "footer", "login")
	//	return
	//}

	//if user.Password == data.Encrypt(r.PostFormValue("password")) {
	//	if user.IsFreelancer(){
	//		if ok := data.CheckFreelancer(user.ID); !ok {
	//			http.Redirect(w, r, "/login", 302)
	//			return
	//		}
	//	} else if user.IsCustomer(){
	//		if ok := data.CheckCustomer(user.ID); !ok {
	//			http.Redirect(w, r, "/login", 302)
	//			return
	//		}
	//	} else {
	//		http.Redirect(w, r, "/login", 302)
	//		return
	//	}
	//	sess, err := user.CreateSession()
	//	if err != nil {
	//		return
	//	}
	//	cookie := http.Cookie{
	//		Name: "_cookie",
	//		Value:    sess.UUID,
	//		HttpOnly: true,
	//		MaxAge: 60 * 60 * 24 * 30,
	//	}
	//	http.SetCookie(w, &cookie)
	//
	//	http.Redirect(w, r, "/my_profile", 302)
	//} else {
	//	pageData.Errors = []string{"Password is wrong"}
	//	//http.Redirect(w, r, "/login", 302)
	//	generateHTML(w, &pageData, nil, "base", "header", "footer", "login")
	//	return
	//}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		sess := data.Session{UUID: cookie.Value}
		err = sess.DeleteByUUID()
		if err != nil{
			log.Println(err)
		}else {
			c := http.Cookie{
				Name:   "_cookie",
				MaxAge: -1}
			http.SetCookie(w, &c)
		}
	} else {
		log.Println(err, " Failed to get cookie")
	}
	http.Redirect(w, r, "/", 302)
}
