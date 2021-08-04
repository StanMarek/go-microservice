package handler

import (
	"context"
	"encoding/json"
	"microservice/src/database"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	result, err := database.DeleteUserById(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(result)
}
