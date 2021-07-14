package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	auth "microservice/src/authentication"
	"microservice/src/database"
	api "microservice/src/httpd/freeAPI"
	endpoint "microservice/src/httpd/handler"

	"github.com/gorilla/mux"
)

func main() {

	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	database.ConnectRedis()

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

	router.HandleFunc("/login", endpoint.Login).Methods("POST")
	router.HandleFunc("/logout", auth.JWTMiddleware(endpoint.Logout)).Methods("POST")
	router.HandleFunc("/register", endpoint.AddUser).Methods("POST") // <- register
	router.HandleFunc("/refreshtoken", auth.RefreshToken).Methods("POST")
	router.HandleFunc("/bored", api.RandomActvity).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}).Methods("GET")
	router.HandleFunc("/users", auth.JWTMiddleware(endpoint.GetAllUsers)).Methods("GET")
	router.HandleFunc("/users/{id}", auth.JWTMiddleware(endpoint.GetUser)).Methods("GET")
	router.HandleFunc("/users", auth.JWTMiddleware(endpoint.AddUser)).Methods("POST") // <- register
	router.HandleFunc("/users/{id}", auth.JWTMiddleware(endpoint.UpdateUser)).Methods("PUT")
	router.HandleFunc("/users/{id}", auth.JWTMiddleware(endpoint.DeleteUser)).Methods("DELETE")

	log.Fatal(server.ListenAndServe())
}
