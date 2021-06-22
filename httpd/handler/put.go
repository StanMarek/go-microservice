package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"microservice/model"

	"github.com/gorilla/mux"
)

type UserUpdateRequest struct {
	Email string `json:"email" validate:"email,required"`
	//Login 		string	`json:"login" validate:"required"`
	Password string `json:"password" validate:"password,required"`
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	indexParam := params["id"]
	indexInt, _ := strconv.Atoi(indexParam)

	var updatedUser UserUpdateRequest
	json.NewDecoder(request.Body).Decode(&updatedUser)

	var user model.User
	var id int
	for index, userIter := range model.Users {
		if userIter.Id == indexInt {
			user = userIter
			id = index
		}
	}

	user.Email = updatedUser.Email
	user.Login = user.ParseEmailToLogin()
	user.Password = updatedUser.Password
	//user.CreatedAt = model.Users[id].CreatedAt
	user.UpdatedAt = time.Now()

	err := user.Validate()
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Error validating user udpdate request: %s", err),
			http.StatusBadRequest,
		)
		return
	} else {
		model.Users[id] = user
		model.Users[id].ToJson(writer)
	}
}
