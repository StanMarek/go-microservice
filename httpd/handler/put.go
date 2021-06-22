package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"microservice/model"

	"github.com/gorilla/mux"
)

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	indexParam := params["id"]
	indexInt, _ := strconv.Atoi(indexParam)
	var updatedUser model.User
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
	for index, user := range model.Users {
		if user.Id == indexInt {
			model.Users[index].Email = updatedUser.Email
			model.Users[index].Login = updatedUser.Login
			model.Users[index].Password = updatedUser.Password
			//json.NewEncoder(writer).Encode(updatedUser)
			updatedUser.ToJson(writer)
		}
	}
}
