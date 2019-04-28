package main

import (
	"graduate/data"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", logging(home))

	// route_auth
	r.HandleFunc("/login", logging(login))
	r.HandleFunc("/registration", logging(registration))
	r.HandleFunc("/registration_account", logging(registrationAccount))
	r.HandleFunc("/login_account", logging(loginAccount))
	r.HandleFunc("/logout", logging(logout))

	// route_user_profile
	subUserProfile := r.PathPrefix("/user").Subrouter()
	subUserProfile.HandleFunc("", logging(freelancerProfile))
	subUserProfile.HandleFunc("/", logging(freelancerProfile))
	subUserProfile.HandleFunc("/about", logging(freelancerProfileAbout))
	subUserProfile.HandleFunc("/works", logging(freelancerProfileWorks))
	subUserProfile.HandleFunc("/contacts", logging(freelancerProfileContacts))
	subUserProfile.HandleFunc("/setting", logging(freelancerProfileSetting))
	// r.HandleFunc("/my_profile", logging(freelancerProfile))
	// r.HandleFunc("/my_profile/about", logging(freelancerProfileAbout))
	// r.HandleFunc("/my_profile/works", logging(freelancerProfileWorks))
	// r.HandleFunc("/my_profile/contacts", logging(freelancerProfileContacts))
	// r.HandleFunc("/my_profile/setting", logging(freelancerProfileSetting))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", r)
	log.Println("Listening...")
	// server := http.Server{
	// 	Addr: "127.0.0.1:8080",
	// }
	// server.ListenAndServe()
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	err := data.SessionChek(r, &session)
	if err != nil {
		data := &Data{"Home", nil}
		generateHTML(w, data, nil, "base", "header", "footer", "home_page")
	} else {
		freelancer, _ := data.GetFreelancerByUserID(session.UserID)
		data := &Data{"Home", &freelancer}
		generateHTML(w, data, nil, "base", "header", "footer", "home_page")
	}
}
