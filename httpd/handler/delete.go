package handler

import (
	"encoding/json"
	"fmt"
	"microservice/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	indexParam := params["id"]
	indexInt, _ := strconv.Atoi(indexParam)
	for index, user := range model.Users {
		if user.Id == indexInt {
			model.Users = append(model.Users[:index], model.Users[index+1:]...)
			fmt.Fprintf(writer, "Deleted user of id: %s", indexParam)
		}
	}
	json.NewEncoder(writer).Encode(model.Users)
}
