package main

import (
	"fmt"
	"graduate/data"
	"html/template"
	"log"
	"net/http"
	"os"
)

func profile(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, err := data.GetFreelancerByUserID(sess.UserID)
		if err != nil {
			log.Println(err)
		}
		pageData := PageData{
			Title :"About",
			User : &user,
		}
		//pageData.Title = "About"
		if data.CheckFreelancer(sess.UserID){
			funcMap, err := freelancerProfile(sess, &pageData)
			if err != nil{
				log.Println(err)
			}
			generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
				"userProfile/freelancer/my_works", "userProfile/freelancer/about")
		}else if data.CheckCustomer(sess.UserID){
			funcMap, err := customerProfile(sess, &pageData)
			if err != nil{
				log.Println(err)
			}
			generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
				"userProfile/customer/my_orders", "userProfile/customer/about")
		}
	}
}

func freelancerProfile(sess data.Session, pageData *PageData) (funcMap template.FuncMap, err error) {
	user, err := data.GetFreelancerByUserID(sess.UserID)
	pageData.User = &user
	if err != nil {
		log.Println(err)
	}
	funcMap = template.FuncMap{
		"getNameSpecialization": data.GetSpecializationName,
	}
	return
}

func customerProfile(sess data.Session, pageData *PageData) (funcMap template.FuncMap, err error) {
	user, err := data.GetCustomerByUserID(sess.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	pageData.User = &user
	pageData.Content = struct {
		Orders []data.Order
	}{user.Orders()}
	return
}

func profileSetting(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title :"Setting",
	}
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user := data.User{}
		if r.Method == http.MethodPost{
			err := r.ParseForm()
			if err != nil {
				panic(err)
			}
			user.FirstName = r.PostFormValue("first_name")
			user.LastName = r.PostFormValue("last_name")
			user.About = r.PostFormValue("about")
			user.Phone = r.PostFormValue("phone")
			user.Facebook = r.PostFormValue("facebook")
			user.Skype = r.PostFormValue("skype")
		}
		if data.CheckFreelancer(sess.UserID){
			settingFreelancer(w, r, sess, &pageData, &user)
		}else if data.CheckCustomer(sess.UserID){
			settingCustomer(w, r, sess, &pageData, &user)
		}
	}
}

func settingFreelancer(w http.ResponseWriter, r *http.Request, sess data.Session, pageData *PageData, user *data.User){
	freelancer, err := data.GetFreelancerByUserID(sess.UserID)
	if err != nil {
		log.Println(err)
	}
	specs, _ := data.GetAllSpecialization()
	if r.Method == http.MethodPost {
		freelancer.User = *user
		freelancer.Specialization = arrayStringToArrayInt(r.Form["specialization[]"])
		err = freelancer.Update()
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/my_profile", 302)
	} else if r.Method == http.MethodGet{
		funcMap := template.FuncMap{
			"getNameSpecialization":  data.GetSpecializationName,
			"containsSpecialization": freelancer.ContainsSpecialization,
		}

		pageData.User = &freelancer
		pageData.Content = struct{
			Specialization []data.Specialization
		}{specs}
		generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
			"userProfile/freelancer/about", "userProfile/setting_base", "userProfile/freelancer/setting",)
	}
}

func settingCustomer(w http.ResponseWriter, r *http.Request, sess data.Session, pageData *PageData, user *data.User) {
	customer, err := data.GetCustomerByUserID(sess.UserID)
	if err != nil {
		log.Println(err)
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		customer.User = *user
		customer.Organization = r.PostFormValue("organization")
		err = customer.Update()
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/my_profile", 302)
	} else {
		pageData.User = &customer
		generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/profile",
			"userProfile/setting_base", "userProfile/customer/setting", "userProfile/customer/about")
	}
}

func newOrder(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}
	customer, err := data.GetCustomerByUserID(sess.UserID)
	if err != nil {
		http.Redirect(w, r, "/my_profile", 302)
		fmt.Println(err)
		return
	}
	pageData := PageData{
		Title :"New order",
		User : &customer,
	}
	if r.Method == http.MethodGet {
		generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/profile",
			"userProfile/customer/about", "userProfile/customer/new_order")
	} else if r.Method == http.MethodPost {
		orderTitle := r.PostFormValue("title")
		orderContent := r.PostFormValue("description")
		customer.CreateOrder(orderTitle, orderContent)
		http.Redirect(w, r,  "/my_profile", 302)
	}
}

func uploadPhoto(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		http.Redirect(w, r,  "/my_profile", 302)
	} else if r.Method == "POST" {
		sess, err := session(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			fmt.Println(err)
			return
		}
		user, _ := data.GetUserByID(sess.UserID)
		path := "static/image/"
		file, handler, err := r.FormFile("photo")
		removePath := user.Photo
		user.Photo, err = uploadFile(path, file, handler)
		if err != nil {
			fmt.Println(err)
			return
		}
		user.UpdatePhoto()
		if removePath != "/static/image/profile.jpg"{
			os.Remove(removePath)
		}
		http.Redirect(w, r,  "/my_profile", 302)
	}
}