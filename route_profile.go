package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"graduate/data"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type PageData struct {
	Title string
	User data.HelperUser
	Content interface{}
	Errors []string
}


func profile(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		pageData := PageData{
			Title :"About",
		}
		if data.CheckFreelancer(sess.UserID){
			funcMap, err := freelancerProfile(sess, &pageData)
			if err != nil{
				log.Println(err)
			}
			generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
				"userProfile/freelancer/my_works", "userProfile/freelancer/about")
		}else if data.CheckCustomer(sess.UserID){
			funcMap, err := customerProfile(sess, &pageData, r)
			if err != nil{
				log.Println(err)
			}
			generateHTML(w, &pageData, funcMap, "base", "header", "footer", "userProfile/profile",
				"userProfile/customer/my_orders", "userProfile/customer/about")
		}
	}
}

func renderDialogsPage(w http.ResponseWriter, r *http.Request){
	pageData := PageData{
		Title :"Dialogs",
	}
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, _ := data.GetUserByID(sess.UserID)
		pageData.User = &user
		if r.Method == http.MethodPost{
			err := r.ParseForm()
			if err != nil {
				panic(err)
			}
		}
		generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/messages")
	}
}

func profileDialogs(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user, _ := data.GetUserByID(sess.UserID)
		dialog, _ := json.Marshal(user.Dialogs())

		w.Header().Set("Content-Type", "application/json")
		w.Write(dialog)
	}
}

func profileDialog(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		vars := mux.Vars(r)
		id, _ := strconv.ParseInt(vars["id"], 10, 8)
		user, _ := data.GetUserByID(sess.UserID)
		dialog, _ := json.Marshal(user.DialogByID(int(id)))

		w.Header().Set("Content-Type", "application/json")
		w.Write(dialog)
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		log.Println(err)
	}
	if r.Method == http.MethodPost{
		//vars := mux.Vars(r)
		//id, _ := strconv.ParseInt(vars["id"], 10, 8)
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		message := &data.Message{}
		jsonErr := json.Unmarshal(body, message)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		log.Println(message)
		user, _ := data.GetUserByID(sess.UserID)
		err = user.SendMessage(message.ReceiverID, message.Text)
		//messageText := r.PostFormValue("message")
		//if err = user.SendMessage(int(id), messageText); err != nil{
		//	log.Println(err)
		//}
	}
}

func profileSetting(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title :"Setting",
	}
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		user := data.User{}
		if r.Method == http.MethodPost{
			err := r.ParseForm()
			if err != nil {
				panic(err)
			}
			user.FirstName = r.PostFormValue("first_name")
			user.LastName = r.PostFormValue("last_name")
			user.About = r.PostFormValue("about")
			user.Phone = r.PostFormValue("phone")
			user.Facebook = r.PostFormValue("facebook")
			user.Skype = r.PostFormValue("skype")
		}
		if data.CheckFreelancer(sess.UserID){
			settingFreelancer(w, r, sess, &pageData, &user)
		}else if data.CheckCustomer(sess.UserID){
			settingCustomer(w, r, sess, &pageData, &user)
		}
	}
}

func uploadPhoto(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		http.Redirect(w, r,  "/my_profile", 302)
	} else if r.Method == "POST" {
		sess, err := session(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			fmt.Println(err)
			return
		}
		user, _ := data.GetUserByID(sess.UserID)
		path := "static/image/"
		file, handler, err := r.FormFile("photo")
		removePath := user.Photo
		user.Photo, err = uploadFile(path, file, handler)
		if err != nil {
			fmt.Println(err)
			return
		}
		user.UpdatePhoto()
		if removePath != "/static/image/profile.jpg"{
			os.Remove(removePath)
		}
		http.Redirect(w, r,  "/my_profile", 302)
	}
}