package main

import (
	"encoding/json"
	"fmt"
	"graduate/data"
	"io/ioutil"
	"log"
	"net/http"
)

func moderatorMain(w http.ResponseWriter, r *http.Request)  {
	http.Redirect(w, r, "/moderator/specializations", 302)
}

func allSpecializations(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, _ := data.GetUserByID(sess.UserID)
	specs, _ := data.GetAllSpecialization()
	pageData := PageData{
		Title:"Specializations",
		User: &user,
		Content: struct {
			Specializations *[]data.Specialization
		}{&specs},
	}
	generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/moderator/panel",
				"userProfile/moderator/specializations")
}

func allUsers(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, _ := data.GetUserByID(sess.UserID)
	users, _ := data.GetAllUsers()
	pageData := PageData{
		Title:"Users",
		User: &user,
		Content: struct {
			Users *[]data.User
		}{&users},
	}
	generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/moderator/panel",
		"userProfile/moderator/users")
}

func allAvailableOrders(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, _ := data.GetUserByID(sess.UserID)
	availableOrders, _ := data.GetOrdersWhere("status_id = $1 ", data.OrderStatusAvailable)
	pageData := PageData{
		Title:"Available orders",
		User: &user,
		Content: struct {
			AvailableOrders *[]data.Order
		}{&availableOrders},
	}
	generateHTML(w, &pageData, nil, "base", "header", "footer", "userProfile/moderator/panel",
		"userProfile/moderator/available_orders")
}

func moderatorCreateSpecialization(w http.ResponseWriter, r *http.Request){
	moderator, err := checkModeratorSession(w, r)
	pageData := PageData{}
	if err != nil {
		log.Println(err)
		return
	}
	if r.Method == http.MethodPost{
		err = r.ParseForm()
		if err != nil {
			pageData.Content = struct {Info string}{"Error parse form"}
			return
		}
		err = moderator.CreateSpecialization(&data.Specialization{
			Name:r.PostFormValue("name"),
		})
		if err != nil{
			pageData.Content = struct {Info string}{"Unsuccessful create specialization: " + err.Error()}
			return
		}
		pageData.Content = struct {Info string}{"Successful create specialization"}
	}
	generateHTML(w, &pageData, nil, "base", "header", "footer", "moderator/new_specialization")
}

func createSpecialization(w http.ResponseWriter, r *http.Request)  {
	sess, err := session(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
		return
	}
	moderator := data.Moderator{}
	moderator.User, _ = data.GetUserByID(sess.UserID)
	if r.Method == http.MethodPost {
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			log.Println(readErr)
		}

		specialization := &data.Specialization{}
		//var name string
		jsonErr := json.Unmarshal(body, specialization)
		if jsonErr != nil {
			log.Println(jsonErr)
		}
		err = moderator.CreateSpecialization(specialization)
		if err != nil{
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		json, _:= json.Marshal(specialization)
		w.Write(json)
	}

	return
}

func moderatorEditSpecialization(w http.ResponseWriter, r *http.Request){
	sess, err := session(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
		return
	}
	moderator := data.Moderator{}
	moderator.User, _ = data.GetUserByID(sess.UserID)
	if r.Method == http.MethodPost {
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			log.Println(readErr)
		}

		specialization := &data.Specialization{}
		//var name string
		jsonErr := json.Unmarshal(body, specialization)
		if jsonErr != nil {
			log.Println(jsonErr)
		}
		err = moderator.UpdateSpecialization(specialization)
		if err != nil{
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		json, _:= json.Marshal(specialization)
		w.Write(json)
	}

	return
}

func moderatorDeleteSpecialization(w http.ResponseWriter, r *http.Request){
	sess, err := session(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
		return
	}
	moderator := data.Moderator{}
	moderator.User, _ = data.GetUserByID(sess.UserID)
	if r.Method == http.MethodPost {
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			log.Println(readErr)
		}

		specialization := &data.Specialization{}
		jsonErr := json.Unmarshal(body, specialization)
		if jsonErr != nil {
			log.Println(jsonErr)
		}
		err = moderator.DeleteSpecialization(specialization)
		if err != nil{
			log.Println(err)
		}
	}

	return
}

func moderatorEditCustomer(w http.ResponseWriter, r *http.Request){

}


func moderatorEditFreelancer(w http.ResponseWriter, r *http.Request){

}

func moderatorDeleteUser(w http.ResponseWriter, r *http.Request){

}

func moderatorEditAvailableOrder(w http.ResponseWriter, r *http.Request){

}

func moderatorDeleteAvailableOrder(w http.ResponseWriter, r *http.Request){

}

func moderatorDeleteRequest(w http.ResponseWriter, r *http.Request){

}

func checkModeratorSession(w http.ResponseWriter, r *http.Request) (moderator data.Moderator, err error) {
	sess, err := session(w, r)
	if err != nil {
		return
	}
	moderator.User, err = data.GetUserByID(sess.ID)
	return
}