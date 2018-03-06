package controllers

import (
	"github.com/gorilla/securecookie"
	"strconv"
	"encoding/json"
	"html/template"
	"net/http"
	"fmt"
	"../models"
)

const tplDir = "./templates"

type Response struct {
	Status string `json:"status"`
	Content string `json:"content"`
}

var index = template.Must(template.ParseFiles(tplDir+"/index.html"))
var header = template.Must(template.ParseFiles(tplDir+"/header.html"))
var footer = template.Must(template.ParseFiles(tplDir+"/footer.html"))
var registro = template.Must(template.ParseFiles(tplDir+"/registro.html"))
var login = template.Must(template.ParseFiles(tplDir+"/login.html"))


func UserRegister(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")
	email := r.FormValue("email")
	w.Header().Set("Content-Type","application/json")
	var res Response
	if username != "" && password != "" && password2 != "" && email != "" && password == password2{
		err := models.RegistrarUsuario(username,password,email)
		if err == nil{
			w.WriteHeader(http.StatusCreated)
			res.Status = "success"
			res.Content = "Registrado con Éxito"
		}else{
			res.Status = "error"
			res.Content = fmt.Sprintf("%s",err)
		}
	} else {
		res.Status = "error"
		res.Content = "Debe completar todos los campos y las contraseñas deben ser iguales."
	}
	json.NewEncoder(w).Encode(res)
}

func UserLogin(w http.ResponseWriter, r *http.Request){
	username := r.FormValue("username")
	password := r.FormValue("password")
	var res Response
	w.Header().Set("Content-Type","application/json")
	if username != "" && password != "" {
		id_user , err := models.Login(username,password)
		if err == nil{
			fmt.Println("usuario logueado: "+strconv.Itoa(id_user))
			setSession(id_user,username,w)
			res.Status = "success"
			res.Content = "Logueado con éxito"
		}else{
			res.Status = "error"
			res.Content = fmt.Sprintf("%s",err)
		}
	} else {
		res.Status = "error"
		res.Content = "Se debe completar todos los campos."
	}
	json.NewEncoder(w).Encode(res)
}

func StaticReg(w http.ResponseWriter, r *http.Request){
	dat := models.PageData{
		Title:		"Inicio - TarinGO",
		Username:	getUserName(r),
	}
	if dat.Username != "" {
		dat.Loged = true
	}
	if dat.Loged {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	header.Execute(w,dat)
	registro.Execute(w,nil)
	footer.Execute(w,nil)
}

func StaticLogin(w http.ResponseWriter, r *http.Request){
	dat := models.PageData{
		Title:		"Inicio - TarinGO",
		Username:	getUserName(r),
	}
	if dat.Username != "" {
		dat.Loged = true
	}
	if dat.Loged {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	header.Execute(w,dat)
	login.Execute(w,nil)
	footer.Execute(w,nil)
}

func Index(w http.ResponseWriter, r *http.Request){
	dat := models.IndexData{}
	dat.Title=		"Inicio - TarinGO"
	dat.Username=	getUserName(r)
	dat.Loged=		false
	if dat.Username != "" {
		dat.Loged = true
	}
	posts , err := getPostsList(0)
	if err != nil {
		fmt.Println(err)
	}
	dat.Posts = posts
	header.Execute(w,dat)
	index.Execute(w,dat)
	footer.Execute(w,nil)
}

func Robots(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, `User-agent: *
					Disallow: /login/
					Disallow: /register/
					Disallow: /reset/
	`)
}

func setSession(id_user int,userName string, response http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
		"user_id": strconv.Itoa(id_user),
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["username"]
		}
	}
	return userName
}