package main

import (
	"encoding/json"
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
	"time"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration
var logger *log.Logger

func init() {
	loadConfig()
	file, err := os.OpenFile("vol.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("-------------------------------------------------------------------")
	logger.Println("Server START ", time.Now().Format("Mon Jan _2 15:04:05 2006"))
	logger.Println("-------------------------------------------------------------------")
}

//generateHTML generate template with passed data and func
func generateHTML(w http.ResponseWriter, data interface{}, funcMap template.FuncMap, filenames ...string) {
	var files []string
	var templates *template.Template

	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates, _ = template.New("").Funcs(funcMap).ParseFiles(files...)
	templates.ExecuteTemplate(w, "base", data)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
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
			err = errors.New("invalid session")
		}
	}
	return
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr + " method: " + r.Method + " - " + r.URL.Path)
		logger.Println(r.RemoteAddr + " method: " + r.Method + " - " + r.URL.Path)
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