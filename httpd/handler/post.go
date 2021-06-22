package handler

import (
	"encoding/json"
	"fmt"
	"microservice/model"
	"net/http"
	"time"
)

type UserPostRequest struct {
	Email string `json:"email" validate:"email,required"`
	//Login 		string	`json:"login" validate:"required"`
	Password string `json:"password" validate:"password,required"`
}

func AddUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")

	var userPostRequest UserPostRequest
	json.NewDecoder(request.Body).Decode(&userPostRequest)

	var user model.User
	user.CreatedAt = time.Now()
	user.Id = model.NextId()
	user.Email = userPostRequest.Email
	user.Login = user.ParseEmailToLogin()
	user.Password = userPostRequest.Password

	user.Validate()
	err := user.Validate()
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Error validating user: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	exists, _ := model.Exists(user.Id)
	if exists {
		fmt.Fprintf(writer, "User of id {%d} already exists", user.Id)
		return
	} else {
		model.Users = append(model.Users, user)
		//json.NewEncoder(writer).Encode(user)
		user.ToJson(writer)
	}
}
