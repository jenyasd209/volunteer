package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"graduate/data"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

func sortWorks(w http.ResponseWriter, r *http.Request) () {
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

func freelancerComment(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}
	pageData := &PageData{
		Title:"Comment order",
	}
	freelancer, _ := data.GetFreelancerByUserID(sess.UserID)
	vars := mux.Vars(r)
	orderID, _ := strconv.ParseInt(vars["id"], 10, 8)
	doneOrder := data.GetDoneOrderByID(int(orderID))
	if r.Method == http.MethodGet {
		pageData.User = &freelancer
		pageData.Content = struct {
			data.CompleteOrder
		}{doneOrder}
		generateHTML(w, &pageData, nil, "base", "header", "footer", "order/freelancer_comment")
	}else if r.Method == http.MethodPost{

		rating, _ := strconv.ParseInt(r.PostFormValue("rating"), 10, 8)
		freelancerComment := data.Comment{
			Text:r.PostFormValue("comment"),
			Rait:float32(rating),
		}
		if doneOrder.FreelancerComment.ID == 0 {
			err = freelancer.CreateComment(&freelancerComment)
			if err != nil{
				log.Println(err)
			}
			doneOrder.FreelancerComment.ID = freelancerComment.ID
			doneOrder.UpdateFreelancerComment()
		}else {
			fmt.Println("Comment exist")
		}
		http.Redirect(w, r, "/my_profile", 302)
	}
}
