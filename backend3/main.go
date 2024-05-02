package main

import (
	"andrey/microservicecontroller"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	port := 5003

	router.HandleFunc("/server1", microservicecontroller.SendPostToService1).Methods("POST")

	fmt.Printf("Server 3 is running at Port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), router))
}
