package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"graduate/data"
	"html/template"
	"net/http"
	"strconv"
)

func viewCustomer(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 8)
	customer, _ := data.GetCustomerByUserID(int(id))
	customerOrders := customer.Orders()
	customerPerformedOrders := customer.PerformedOrders()
	customerDoneOrders := customer.CompleteOrders()
	pageData := PageData{
		Title :customer.FirstName + " " + customer.LastName,
		Content: struct {
			*data.Customer
			Orders *[]data.Order
			PerformedOrders *[]data.PerformedOrder
			DoneOrders *[]data.CompleteOrder

		}{&customer, &customerOrders, &customerPerformedOrders,
			&customerDoneOrders},
	}
	sess, err := session(w, r)
	if err == nil {
		user, _ := data.GetUserByID(sess.UserID)
		pageData.User = &user
	}
	funcMap := template.FuncMap{
		"getNameSpecialization":  data.GetSpecializationName,
	}
	generateHTML(w, &pageData, funcMap, "base", "header", "footer", "customer/customer_view")
}

func customerSortOrders(w http.ResponseWriter, r *http.Request,) () {
	vars := mux.Vars(r)
	status := vars["status"]
	customerID, _ := strconv.ParseInt(vars["id"], 10, 8)
	var orders []byte

	customer, _ := data.GetCustomerByUserID(int(customerID))
	if status == "" {
		orders, _ = json.Marshal(customer.Orders())
	}else if status == "available" {
		orders, _ = json.Marshal(customer.GetOrdersByStatus(data.OrderStatusAvailable))
	}else if status == "done" {
		orders, _ = json.Marshal(customer.CompleteOrders())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(orders)
	return
}