package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

// func session(w http.ResponseWriter, r *http.Request, session *user.Sessionable) (err error) {
// 	cookie, err := r.Cookie("_cookie")
// 	if err == nil {
// 		// sess = session
// 		(*session).SetUUID(cookie.Value)
// 		// session := freelancer.Session{UUID: cookie.Value}
// 		if ok, err := (*session).Check(); !ok {
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
