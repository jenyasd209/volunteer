package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"graduate/data"
	"log"
	"net/http"
	"strconv"
)

func allOrders(w http.ResponseWriter, r *http.Request)  {
	pageData := PageData{
		Title:"Jobs",
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
	orders := new([]data.Order)
	specialization, _ := data.GetAllSpecialization()
	if search := r.FormValue("search"); search != ""{
		*orders, err = data.GetOrdersWhere(`title ILIKE '%' || $1 || '%'`, search)
		if err != nil{
			log.Println(err)
		}
		if len(*orders) == 0{
			orders = nil
		}
	}else{
		*orders, err = data.GetAllOrders()
		if err != nil{
			log.Println(err)
		}
		if len(*orders) == 0{
			orders = nil
		}
	}
	pageData.Content = struct {
		Orders *[]data.Order
		Specialization *[]data.Specialization
	}{orders, &specialization}

	generateHTML(w, &pageData, nil, "base", "header", "footer", "order/orders")
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
			"userProfile/customer/about", "order/new_order")
	} else if r.Method == http.MethodPost {
		orderTitle := r.PostFormValue("title")
		orderContent := r.PostFormValue("description")
		customer.CreateOrder(orderTitle, orderContent)
		http.Redirect(w, r,  "/my_profile", 302)
	}
}

func viewOrder(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	order := data.GetOrderByID(int(id))
	if (data.Order{}) == order{
		http.Redirect(w, r, "/not_found", 302)
	}
	order.FreelancerRequest = order.GetRequests()
	customer, _ := data.GetCustomerByUserID(order.Customer.User.ID)
	pageData := PageData{
		Title :order.Title,
		Content: struct {
			*data.Customer
			*data.Order
		}{&customer, &order},
	}
	sess, err := session(w, r)
	if err == nil {
		user, _ := data.GetUserByID(sess.UserID)
		pageData.User = &user
	}
	generateHTML(w, &pageData, nil, "base", "header", "footer", "order/order_view")
}

func editOrder(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	order := data.GetOrderByID(int(id))
	customer, err := data.GetCustomerByUserID(sess.UserID)
	if err != nil {
		http.Redirect(w, r, "/my_profile", 302)
		fmt.Println(err)
		return
	}
	pageData := PageData{
		Title :"Edit order \"" + order.Title + "\"",
		User : &customer,
		Content: struct {
			data.Order
		}{order},
	}
	if r.Method == http.MethodGet {
		generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/profile",
			"userProfile/customer/about", "order/edit_order")
	} else if r.Method == http.MethodPost {
		order.Title = r.PostFormValue("title")
		order.Content = r.PostFormValue("description")
		err := customer.UpdateOrder(&order)
		if err != nil{
			log.Println(err)
		}
		http.Redirect(w, r,  "/my_profile", 302)
	}
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	order := data.GetOrderByID(int(id))
	customer, err := data.GetCustomerByUserID(sess.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = customer.DeleteOrder(order)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/my_profile", 302)
}

func newRequest(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	vars := mux.Vars(r)
	orderID, _ := strconv.ParseInt(vars["id"], 10, 8)
	if r.Method == http.MethodPost{
		freelancer, err := data.GetFreelancerByUserID(sess.UserID)
		if err != nil {
			log.Println(err)
			return
		}
		addText := r.PostFormValue("text")
		err = freelancer.CreateRequest(int(orderID), addText)
		if err != nil {
			log.Println(err)
		}
	}
	http.Redirect(w, r, "/orders/id" + vars["id"], 302)
}