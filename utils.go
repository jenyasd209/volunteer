package main

import (
	"errors"
	"fmt"
	"graduate/data"
	"html/template"
	"io"
	"log"
	"mime/multipart"
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

//var pageData PageData

func generateHTML(w http.ResponseWriter, data interface{}, funcMap template.FuncMap, filenames ...string) {
	var files []string
	var templates *template.Template

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates, _ = template.New("").Funcs(funcMap).ParseFiles(files...)
	templates.ExecuteTemplate(w, "base", data)
	//pageData.Content = nil
	//pageData.Errors = nil
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

func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = data.Session{UUID: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session.")
		}
	}
	return
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr + " method: " + r.Method + " - " + r.URL.Path)
		f(w, r)
	}
}

func uploadFile(path string, file multipart.File, handler *multipart.FileHeader) (filepath string, err error){
	defer file.Close()
	f, err := os.OpenFile(path + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return
	}
	filepath = "/" + path + handler.Filename
	return
}
