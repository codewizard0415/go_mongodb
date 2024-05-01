package main

import (
	"log"
	"net/http"

	"mongo/controllers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	s.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	s.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")
	s.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	s.HandleFunc("/users", controllers.GetUsers).Methods("GET")

	log.Fatal(http.ListenAndServe(":9001", s))
}
