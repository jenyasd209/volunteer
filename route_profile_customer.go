package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"graduate/data"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

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
	funcMap = template.FuncMap{
		"setStars": setRaiting,
	}
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
		funcMap := template.FuncMap{
			"setStars": setRaiting,
		}
		generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
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
		funcMap := template.FuncMap{
			"setStars": setRaiting,
		}
		generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
			"userProfile/customer/about", "userProfile/customer/new_order")
	} else if r.Method == http.MethodPost {
		orderTitle := r.PostFormValue("title")
		orderContent := r.PostFormValue("description")
		customer.CreateOrder(orderTitle, orderContent)
		http.Redirect(w, r,  "/my_profile", 302)
	}
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
		funcMap := template.FuncMap{
			"setStars": setRaiting,
		}
		generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
			"userProfile/customer/about", "userProfile/customer/edit_order")
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

func sortOrder(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		pageData := PageData{
			Title :"About",
		}
		if data.CheckFreelancer(sess.UserID){
			user, err := data.GetCustomerByUserID(sess.UserID)
			if err != nil {
				log.Println(err)
				return
			}
			vars := mux.Vars(r)
			pageData.User = &user
			orders := user.Orders()
			if vars["param"] == "title" {
				sort.Slice(orders, func(i, j int) bool { return orders[i].Title < orders[j].Title })
			}
			if vars["param"] == "date" {
				http.Redirect(w, r, "/my_profile", 302)
			}
			if vars["param"] == "available" {
				user.GetOrdersByStatus(data.OrderStatusAvailable)
			}
			if vars["param"] == "performed" {
				user.GetOrdersByStatus(data.OrderStatusPerformed)
			}
			if vars["param"] == "done" {
				user.GetOrdersByStatus(data.OrderStatusDone)
			}
			pageData.Content = struct {
				Orders *[]data.Order
			}{&orders}

			funcMap := template.FuncMap{
				"setStars": setRaiting,
			}
			generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
				"userProfile/freelancer/my_works", "userProfile/freelancer/about")
		}
	}
}