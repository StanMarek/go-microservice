package handler

import (
	"encoding/json"
	"microservice/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	json.NewEncoder(writer).Encode(model.Users)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	indexParam := params["id"]
	indexInt, _ := strconv.Atoi(indexParam)
	var jsonUser model.User
	for _, user := range model.Users {
		if user.Id == indexInt {
			jsonUser = user
			break
		}
	}
	//json.NewEncoder(writer).Encode(jsonUser)
	jsonUser.ToJson(writer)
}
