package main

import (
	"github.com/gorilla/mux"
	"graduate/data"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func allFreelancers(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title: "Freelancers",
	}
	sess, err := session(w, r)
	if err == nil {
		user, _ := data.GetUserByID(sess.UserID)
		if user.IsFreelancer() {
			user, _ := data.GetFreelancerByUserID(user.ID)
			pageData.User = &user
		} else if user.IsCustomer() {
			user, _ := data.GetCustomerByUserID(user.ID)
			pageData.User = &user
		}else {
			user, _ := data.GetUserByID(user.ID)
			pageData.User = &user
		}
	}
	freelancers := new([]data.Freelancer)
	specialization, _ := data.GetAllSpecialization()
	if search := r.FormValue("search"); search != "" {
		*freelancers, err = data.GetFreelancersWhere(`first_name ILIKE '%' || $1 || '%'
													  OR last_name ILIKE '%' || $1 || '%'`, search)
		if err != nil {
			log.Println(err)
		}
		if len(*freelancers) == 0 {
			freelancers = nil
		}
	} else {
		*freelancers, err = data.GetAllFreelancers()
		if err != nil {
			log.Println(err)
		}
		if len(*freelancers) == 0 {
			freelancers = nil
		}
	}
	pageData.Content = struct {
		Freelancers     *[]data.Freelancer
		Specializations *[]data.Specialization
	}{freelancers, &specialization}

	generateHTML(w, &pageData, nil, "base", "header", "footer", "freelancer/freelancers")
}

func viewFreelancer(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	freelancer, _ := data.GetFreelancerByUserID(int(id))
	freelancerOrders := freelancer.FinishWorks()
	pageData := PageData{
		Title :freelancer.FirstName + " " + freelancer.LastName,
		Content: struct {
			*data.Freelancer
			FreelancerOrders *[]data.CompleteOrder
		}{&freelancer, &freelancerOrders},
	}
	sess, err := session(w, r)
	if err == nil {
		user, _ := data.GetUserByID(sess.UserID)
		pageData.User = &user
	}
	funcMap := template.FuncMap{
		"getNameSpecialization":  data.GetSpecializationName,
	}
	generateHTML(w, &pageData, funcMap, "base", "header", "footer", "freelancer/freelancer_view")
}

func specialization(w http.ResponseWriter, r *http.Request)  {
	pageData := PageData{
		Title:"Freelancers",
	}
	sess, err := session(w, r)
	if err == nil {
		user, _ := data.GetUserByID(sess.UserID)
		if user.IsFreelancer(){
			user, _ := data.GetFreelancerByUserID(user.ID)
			pageData.User = &user
		}else if user.IsCustomer(){
			user, _ := data.GetCustomerByUserID(user.ID)
			pageData.User = &user
		}
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	freelancers := new([]data.Freelancer)
	specializations, _ := data.GetAllSpecialization()
	if search := r.FormValue("search"); search != ""{
		*freelancers, err = data.GetFreelancersWhere(`first_name ILIKE '%' || $1 || '%'
													  OR last_name ILIKE '%' || $1 || '%'`, search)
		if err != nil{
			log.Println(err)
		}
		if len(*freelancers) == 0{
			freelancers = nil
		}
	}else{
		*freelancers, err = data.GetFreelancersWhere(" $1 = ANY (specialization)", id)
		if err != nil{
			log.Println(err)
		}
		if len(*freelancers) == 0{
			freelancers = nil
		}
	}

	pageData.Content = struct {
		Freelancers *[]data.Freelancer
		Specializations *[]data.Specialization
	}{freelancers, &specializations}
	generateHTML(w, &pageData, nil, "base", "header", "footer", "freelancer/freelancers")
}