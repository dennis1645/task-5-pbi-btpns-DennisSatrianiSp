package main

import (
	"log"
	"net/http"

	"github.com/dennis1645/go-api/controllers/authcontroller"
	"github.com/dennis1645/go-api/controllers/dashboardcontroller"
	"github.com/dennis1645/go-api/controllers/photoscontroller"
	"github.com/dennis1645/go-api/controllers/userscontroller"
	"github.com/dennis1645/go-api/middleware"
	"github.com/dennis1645/go-api/models"
	"github.com/gorilla/mux"
)

func main() {

	models.ConnectDatabase()

	root := mux.NewRouter()

	root.HandleFunc("/users/login", authcontroller.Login).Methods("POST")
	root.HandleFunc("/users/register", authcontroller.Register).Methods("POST")
	root.HandleFunc("/users/logout", authcontroller.Logout).Methods("GET")

	api := root.PathPrefix("/api").Subrouter()
	api.HandleFunc("/dashboard/{id}", dashboardcontroller.Index).Methods("GET")
	api.HandleFunc("/photos/{id}", photoscontroller.Show).Methods("GET")
	api.HandleFunc("/add/{id}", photoscontroller.Create).Methods("POST")
	api.HandleFunc("/update/{userId}/{id}", photoscontroller.Update).Methods("PUT")
	api.HandleFunc("/delete/{userId}/{id}", photoscontroller.Delete).Methods("DELETE")
	api.HandleFunc("/users/update/{userId}", userscontroller.Update).Methods("PUT")
	api.HandleFunc("/users/delete/{userId}", userscontroller.Delete).Methods("DELETE")
	api.Use(middleware.JwtMiddleware)

	log.Fatal(http.ListenAndServe(":8080", root))
}
