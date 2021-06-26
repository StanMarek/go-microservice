package handler

import (
	"encoding/json"
	"log"
	"microservice/database"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")

	users, err := database.GetUsers()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(writer).Encode(users)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")

	params := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	user, err := database.GetUserById(id)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	user.ToJson(writer)
}
