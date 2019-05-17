package main

import (
	"github.com/gorilla/mux"
	"graduate/data"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.HandleFunc("/not_found", logging(notFound))

	r.HandleFunc("/", logging(home))

	// route_auth
	r.HandleFunc("/login", logging(login))
	r.HandleFunc("/registration", logging(registration))
	r.HandleFunc("/registration_account", logging(registrationAccount))
	r.HandleFunc("/login_account", logging(loginAccount))
	r.HandleFunc("/logout", logging(logout))

	// route_user_profile
	subUserProfile := r.PathPrefix("/my_profile").Subrouter()
	subUserProfile.HandleFunc("", logging(profile))
	subUserProfile.HandleFunc("/", logging(profile))
	subUserProfile.HandleFunc("/setting", logging(profileSetting))
	subUserProfile.HandleFunc("/upload_photo", logging(uploadPhoto))

	//route customer profile
	//subUserProfile.HandleFunc("/sort_by", logging(sortOrder))
	subUserProfile.HandleFunc("/new_order", logging(newOrder))
	subUserProfile.HandleFunc("/edit_order/id{id:[0-9]+}", logging(editOrder))
	subUserProfile.HandleFunc("/delete_order/id{id:[0-9]+}", logging(deleteOrder))

	//route freelancer profile

	//route freelancers
	subFreelancers := r.PathPrefix("/freelancers").Subrouter()
	subFreelancers.HandleFunc("", logging(allFreelancers))
	subFreelancers.HandleFunc("/", logging(allFreelancers))
	subFreelancers.HandleFunc("/id{id:[0-9]+}", logging(viewFreelancer))

	//route orders
	subOrder := r.PathPrefix("/orders").Subrouter()
	subOrder.HandleFunc("", logging(allOrders))
	subOrder.HandleFunc("/", logging(allOrders))
	subOrder.HandleFunc("/id{id:[0-9]+}", logging(viewOrder))
	subOrder.HandleFunc("/id{id:[0-9]+}/new_request", logging(newRequest))
	subOrder.HandleFunc("/search", logging(allOrders))


	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", r)
	log.Println("Listening...")
	server := http.Server{
		Addr: config.Address,
	}
	server.ListenAndServe()
	//http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	pageData := &PageData{
		Title : "Home",
	}
	sess, err := session(w, r)
	if err == nil {
		user, _ := data.GetUserByID(sess.UserID)
		if user.IsFreelancer(){
			freelancer, _ := data.GetFreelancerByUserID(sess.UserID)
			pageData.User = &freelancer
		}else if user.IsCustomer(){
			customer, _ := data.GetCustomerByUserID(sess.UserID)
			pageData.User = &customer
		}
	}
	generateHTML(w, pageData, nil, "base", "header", "footer", "home_page")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	pageData := &PageData{
		Title : "404 Not Found",
	}
	sess, err := session(w, r)
	if err == nil {
		user, _ := data.GetUserByID(sess.UserID)
		if user.IsFreelancer(){
			freelancer, _ := data.GetFreelancerByUserID(sess.UserID)
			pageData.User = &freelancer
		}else if user.IsCustomer(){
			customer, _ := data.GetCustomerByUserID(sess.UserID)
			pageData.User = &customer
		}
	}
	generateHTML(w, pageData, nil, "not_found")
}
