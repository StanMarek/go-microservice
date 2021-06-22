package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	endpoint "microservice/httpd/handler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9090",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	HandleRequest(server, router)
}

func HandleRequest(server *http.Server, router *mux.Router) {

	router.HandleFunc("/root", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}).Methods("GET")
	router.HandleFunc("/users", endpoint.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", endpoint.GetUser).Methods("GET")
	router.HandleFunc("/users", endpoint.AddUser).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", endpoint.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id:[0-9]+}", endpoint.DeleteUser).Methods("DELETE")

	log.Fatal(server.ListenAndServe())
}