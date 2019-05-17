package main

import (
	"fmt"
	"graduate/data"
	"html/template"
	"log"
	"net/http"
)

func customerProfile(sess data.Session, pageData *PageData, r *http.Request) (funcMap template.FuncMap, err error) {
	user, err := data.GetCustomerByUserID(sess.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	pageData.User = &user
	//orders := new([]data.Order)
	//if param := r.Form.Get("param"); param != ""{
	//	vars := mux.Vars(r)
	//	orders = sortOrder(vars["param"], user)
	//	log.Println("param")
	//}else {
	//	log.Println("no param")
	//	*orders = user.Orders()
	//}
	pageData.Content = struct {
		Orders []data.Order
	}{user.Orders()}
	return
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
		user.ID = customer.User.ID
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

func sortOrder(param string, user data.Customer) (orders *[]data.Order) {
	if param == "all" {
		*orders = user.Orders()
	}else if param == "available" {
		orders = user.GetOrdersByStatus(data.OrderStatusAvailable)
	}else if param == "performed" {
		orders = user.GetOrdersByStatus(data.OrderStatusPerformed)
	}else if param == "done" {
		orders = user.GetOrdersByStatus(data.OrderStatusDone)
	}
	return
}