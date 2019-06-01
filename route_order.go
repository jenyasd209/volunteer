package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"graduate/data"
	"html/template"
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
		*orders, err = data.GetOrdersWhere(`title ILIKE '%' || $1 || '%' and status_id = 1`, search)
		if err != nil{
			log.Println(err)
		}
		if len(*orders) == 0{
			orders = nil
		}
	}else{
		//*orders, err = data.GetAllOrders()
		*orders, err = data.GetOrdersWhere(`title ILIKE '%' || $1 || '%' and status_id = 1`, search)
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
	specs, _ := data.GetAllSpecialization()
	pageData := PageData{
		Title :"New order",
		User : &customer,
		Content: struct {
			Specialization *[]data.Specialization
		}{&specs},
	}
	if r.Method == http.MethodGet {
		generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/profile",
			"userProfile/customer/about", "order/new_order")
	} else if r.Method == http.MethodPost {
		specID, _ := strconv.ParseInt(r.PostFormValue("specialization"), 10, 8)
		newOrder := &data.Order{
			Title: r.PostFormValue("title"),
			Content: r.PostFormValue("description"),
			SpecializationID: int(specID),
		}
		customer.CreateOrder(newOrder)
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
	funcMap := new(template.FuncMap)
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
	funcMap = &template.FuncMap{
		"existFreelanserOrders":  data.ExistOffer,
	}
	generateHTML(w, &pageData, *funcMap, "base", "header", "footer", "order/order_view")
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
	specs, _ := data.GetAllSpecialization()
	specID, _ := strconv.ParseInt(r.PostFormValue("specialization"), 10, 8)
	pageData := PageData{
		Title :"Edit order \"" + order.Title + "\"",
		User : &customer,
		Content: struct {
			data.Order
			Specialization *[]data.Specialization
		}{order, &specs},
	}
	if r.Method == http.MethodGet {
		generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/profile",
			"userProfile/customer/about", "order/edit_order")
	} else if r.Method == http.MethodPost {
		order.Title = r.PostFormValue("title")
		order.Content = r.PostFormValue("description")
		order.SpecializationID = int(specID)
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

func categoryOrders(w http.ResponseWriter, r *http.Request)  {
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
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	//if err != nil {
	//	http.Redirect(w, r, "/my_profile", 302)
	//	fmt.Println(err)
	//	return
	//}
	orders := new([]data.Order)
	specializations, _ := data.GetAllSpecialization()
	if search := r.FormValue("search"); search != ""{
		*orders, err = data.GetOrdersWhere(`title ILIKE '%' || $1 || '%' and status_id = 1`, search)
		if err != nil{
			log.Println(err)
		}
		if len(*orders) == 0{
			orders = nil
		}
	}else{
		*orders, err = data.GetOrdersWhere(" specialization_id = $1 and status_id = 1", id)
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
	}{orders, &specializations}
	generateHTML(w, &pageData, nil, "base", "header", "footer", "order/orders")
}

func existFreelanserOrders(freelancerID int) (result bool) {
	return
}