package handler

import (
	"encoding/json"
	"fmt"
	"microservice/model"
	"net/http"
)

func AddUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	var user model.User
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
