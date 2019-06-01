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

func customerProfile(sess data.Session, pageData *PageData, r *http.Request) (funcMap template.FuncMap, err error) {
	user, err := data.GetCustomerByUserID(sess.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	pageData.User = &user
	//orders := new([]data.CardOrder)
	//if param := r.Form.Get("param"); param != ""{
	//	vars := mux.Vars(r)
	//	orders = sortOrder(vars["param"], user)
	//	log.Println("param")
	//}else {
	//	log.Println("no param")
	//	*orders = user.Orders()
	//}
	specs, _ := data.GetAllSpecialization()
	pageData.Content = struct {
		Orders []data.Order
		Specialization *[]data.Specialization
	}{user.Orders(), &specs}
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

func sortOrders(w http.ResponseWriter, r *http.Request,) () {
	vars := mux.Vars(r)
	status := vars["status"]
	var orders []byte
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}

	customer, _ := data.GetCustomerByUserID(sess.UserID)
	if status == "" {
		orders, _ = json.Marshal(customer.Orders())
	}else if status == "available" {
		orders, _ = json.Marshal(customer.GetOrdersByStatus(data.OrderStatusAvailable))
	}else if status == "performed" {
		orders, _ = json.Marshal(customer.PerformedOrders())
	}else if status == "done" {
		orders, _ = json.Marshal(customer.CompleteOrders())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(orders)
	return
}

func selectRequest(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}
	customer, _ := data.GetCustomerByUserID(sess.UserID)
	vars := mux.Vars(r)
	orderID, _ := strconv.ParseInt(vars["id"], 10, 8)
	freelancerID, _ := strconv.ParseInt(vars["freelancer_id"], 10, 8)

	freelancer, _ := data.GetFreelancerByUserID(int(freelancerID))

	performerOrder := &data.PerformedOrder{
		Order: data.GetOrderByID(int(orderID)),
		Freelancer: freelancer,
	}

	err = customer.MakeOrderPerformed(performerOrder)
	if err != nil{
		log.Println(err)
	}
	http.Redirect(w, r, "/orders/id"+vars["id"], 302)
}

func orderDone(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		fmt.Println(err)
		return
	}
	pageData := &PageData{
		Title:"Make order done",
	}
	customer, _ := data.GetCustomerByUserID(sess.UserID)
	vars := mux.Vars(r)
	orderID, _ := strconv.ParseInt(vars["id"], 10, 8)
	performedOrder := data.GetPerformedOrdersByID(int(orderID))
	if r.Method == http.MethodGet {
		pageData.User = &customer
		pageData.Content = struct {
			data.PerformedOrder
		}{performedOrder}
		generateHTML(w, &pageData, nil, "base", "header", "footer", "order/order_make_done")
	}else if r.Method == http.MethodPost{

		rating, _ := strconv.ParseInt(r.PostFormValue("rating"), 10, 8)
		freelancer, _ := data.GetFreelancerByUserID(performedOrder.Freelancer.User.ID)
		customerComment := data.Comment{
			Text:r.PostFormValue("comment"),
			Rait:float32(rating),
		}
		doneOrder := &data.CompleteOrder{
			Order: data.GetOrderByID(int(orderID)),
			Freelancer: freelancer,
			CustomerComment:customerComment,
		}
		err = customer.MakeOrderDone(doneOrder)
		if err != nil{
			log.Println(err)
		}
		http.Redirect(w, r, "/my_profile", 302)
	}
	//vars := mux.Vars(r)
	//orderID, _ := strconv.ParseInt(vars["id"], 10, 8)
	//freelancerID, _ := strconv.ParseInt(vars["freelancer_id"], 10, 8)
	//
	//freelancer, _ := data.GetFreelancerByUserID(int(freelancerID))
	//doneOrder := &data.CompleteOrder{
	//	Order: data.GetOrderByID(int(orderID)),
	//	Freelancer: freelancer,
	//}
	//err = customer.MakeOrderDone(doneOrder)
	//if err != nil{
	//	log.Println(err)
	//}
}