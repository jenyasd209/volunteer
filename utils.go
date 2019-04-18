package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Data struct {
	PageTitle string
	Other     interface{}
}

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "base", data)
}

//arrayStringToArrayInt conver []string to []int
func arrayStringToArrayInt(strArray []string) (intArry []int) {
	for _, i := range strArray {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		intArry = append(intArry, j)
	}
	return
}

// func checkAuth(f http.HandleFunc) http.HandleFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
//
// 		f(w http.ResponseWriter, r *http.Request)
// 	}
// }

// func sess(w http.ResponseWriter, r *http.Request) (err error) {
// 	cookie, err := r.Cookie("_cookie")
// 	if err == nil {
// 		sess := data.Session{}
// 		sess.SetUUID(cookie.Value)
// 		// session := freelancer.Session{UUID: cookie.Value}
// 		if ok, err := session.Check(); !ok {
// 			fmt.Println(err)
// 		}
// 	}
//
// 	return
// }

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr + " - " + r.URL.Path)
		f(w, r)
	}
}
