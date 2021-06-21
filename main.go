package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

// TODO: add more validation rules
type User struct {
	Id       int    `json:"_id,omitempty"`
	Email    string `json:"email,omitempty" validate:"email,required"`
	Login    string `json:"login,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"password,required"`
}

func (u *User) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}

// TODO: correct password validation
func PasswordValidation(fl validator.FieldLevel) bool {
	const minLength = 6
	var upperCase bool = false
	var lowerCase bool = false
	var number bool = false
	var currentLength = 0
	password := fl.Field().String()

	for _, character := range password {
		if unicode.IsNumber(character) {
			number = true
			currentLength++
		}
		if unicode.IsUpper(character) {
			upperCase = true
			currentLength++
		}
		if unicode.IsLower(character) {
			lowerCase = true
			currentLength++
		}
	}
	if upperCase || lowerCase || number || currentLength >= minLength {
		return true
	} else {
		return false
	}
}

func (u *User) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", PasswordValidation)
	return validate.Struct(u)
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
	//json.NewEncoder(writer).Encode(jsonUser)
	jsonUser.ToJson(writer)
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
	err := updatedUser.Validate()
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Error validating user: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	for index, user := range Users {
		if user.Id == indexInt {
			Users[index].Email = updatedUser.Email
			Users[index].Login = updatedUser.Login
			Users[index].Password = updatedUser.Password
			//json.NewEncoder(writer).Encode(updatedUser)
			updatedUser.ToJson(writer)
		}
	}
}

func AddUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	//user.Validate()
	err := user.Validate()
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Error validating user: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	exists, _ := Exists(user.Id)
	if exists {
		fmt.Fprintf(writer, "User of id {%d} already exists", user.Id)
		return
	} else {
		Users = append(Users, user)
		//json.NewEncoder(writer).Encode(user)
		user.ToJson(writer)
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
