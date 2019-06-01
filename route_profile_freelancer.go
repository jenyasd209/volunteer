package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"graduate/data"
	"html/template"
	"log"
	"net/http"
)

func freelancerProfile(sess data.Session, pageData *PageData) (funcMap template.FuncMap, err error) {
	user, err := data.GetFreelancerByUserID(sess.UserID)
	pageData.User = &user
	if err != nil {
		log.Println(err)
	}
	freelancerCompleteOrders := user.FinishWorks()
	pageData.Content = struct {
		FreelancerCompleteOrders *[]data.CompleteOrder
	}{&freelancerCompleteOrders}

	funcMap = template.FuncMap{
		"getNameSpecialization": data.GetSpecializationName,
	}
	return
}

func settingFreelancer(w http.ResponseWriter, r *http.Request, sess data.Session, pageData *PageData, user *data.User){
	freelancer, err := data.GetFreelancerByUserID(sess.UserID)
	if err != nil {
		log.Println(err)
	}
	specs, _ := data.GetAllSpecialization()
	if r.Method == http.MethodPost {
		user.ID = freelancer.User.ID
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

func sortWorks(w http.ResponseWriter, r *http.Request,) () {
	vars := mux.Vars(r)
	status := vars["status"]
	var orders []byte
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}

	freelancer, _ := data.GetFreelancerByUserID(sess.UserID)
	if status == "" {
		orders, _ = json.Marshal(freelancer.FinishWorks())
	}else if status == "performed" {
		orders, _ = json.Marshal(freelancer.PerformingOrders())
	}else if status == "done" {
		orders, _ = json.Marshal(freelancer.FinishWorks())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(orders)
	return
}
