package main

import (
	"fmt"
	"graduate/data"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageData struct {
	Title string
	User data.HelperUser
	Content interface{}
	Error string
}

//type PageData struct {
//	PageTitle string
//	Content struct{
//		Data interface{}
//		User data.HelperUser
//	}
//}

//var pageData PageData

func generateHTML(w http.ResponseWriter, data interface{}, funcMap template.FuncMap, filenames ...string) {
	var files []string
	var templates *template.Template

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	// templates = template.Must(templates.Funcs(funcMap).ParseFiles(files...))
	templates, _ = template.New("").Funcs(funcMap).ParseFiles(files...)

	templates.ExecuteTemplate(w, "base", data)
}

//arrayStringToArrayInt convert []string to []int
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
