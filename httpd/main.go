package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"microservice/authentication"
	"microservice/database"
	endpoint "microservice/httpd/handler"

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
	router.HandleFunc("/logout", JWTMiddleware(endpoint.Logout)).Methods("POST")
	router.HandleFunc("/register", endpoint.AddUser).Methods("POST") // <- register

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}).Methods("GET")
	router.HandleFunc("/users", JWTMiddleware(endpoint.GetAllUsers)).Methods("GET")
	router.HandleFunc("/users/{id}", JWTMiddleware(endpoint.GetUser)).Methods("GET")
	router.HandleFunc("/users", JWTMiddleware(endpoint.AddUser)).Methods("POST") // <- register
	router.HandleFunc("/users/{id}", JWTMiddleware(endpoint.UpdateUser)).Methods("PUT")
	router.HandleFunc("/users/{id}", JWTMiddleware(endpoint.DeleteUser)).Methods("DELETE")

	log.Fatal(server.ListenAndServe())
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isTokenValidError := authentication.IsTokenValid(request)
		if isTokenValidError != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"message": ` + isTokenValidError.Error() + `"}`))
			return
		}
		next.ServeHTTP(writer, request)
	})
}
