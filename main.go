package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"./controllers"
)

var router = mux.NewRouter()

func main() {
	// Routes
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/robots.txt", controllers.Robots).Methods("GET")
	router.HandleFunc("/registro", controllers.UserRegister).Methods("POST")
	router.HandleFunc("/registro", controllers.StaticReg).Methods("GET")
	router.HandleFunc("/login", controllers.StaticLogin).Methods("GET")
	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")

	router.HandleFunc("/", controllers.Index).Methods("GET")

	log.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080",router)
	if (err != nil){
		panic(err)
	}
}