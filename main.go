package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Id       int    `json:"_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

var Users []User

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/root", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}).Methods("GET")
	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}", GetUser).Methods("GET")
	router.HandleFunc("/user", AddUser).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}", UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id:[0-9]+}", DeleteUser).Methods("DELETE")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9090",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	indexParam := params["id"]
	indexInt, _ := strconv.Atoi(indexParam)
	var jsonUser User
	for _, user := range Users {
		if user.Id == indexInt {
			jsonUser = user
			break
		}
	}
	json.NewEncoder(writer).Encode(jsonUser)
}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	indexParam := params["id"]
	indexInt, _ := strconv.Atoi(indexParam)
	for index, user := range Users {
		if user.Id == indexInt {
			Users = append(Users[:index], Users[index+1])
			fmt.Fprintf(writer, "Deleted user of id: %s", indexParam)
		}
	}
	json.NewEncoder(writer).Encode(Users)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	indexParam := params["id"]
	indexInt, _ := strconv.Atoi(indexParam)
	var updatedUser User
	json.NewDecoder(request.Body).Decode(&updatedUser)
	for index, user := range Users {
		if user.Id == indexInt {
			Users[index].Email = updatedUser.Email
			Users[index].Login = updatedUser.Login
			Users[index].Password = updatedUser.Password
			json.NewEncoder(writer).Encode(updatedUser)
		}
	}
}

func AddUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	exists, _ := Exists(user.Id)
	if exists {
		fmt.Fprintf(writer, "User od id %d already exists", user.Id)
		return
	} else {
		Users = append(Users, user)
		json.NewEncoder(writer).Encode(user)
	}
}

func GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	json.NewEncoder(writer).Encode(Users)
}

func Exists(id int) (bool, int) {
	for index, i := range Users {
		if id == i.Id {
			return true, index
		}
	}
	return false, -1
}
