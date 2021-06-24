package handler

import (
	"context"
	"encoding/json"
	"microservice/model"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)
	collection := client.Database("go_microservice").Collection("users")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// indexParam := params["id"]
	// indexInt, _ := strconv.Atoi(indexParam)
	// for index, user := range model.Users {
	// 	if user.Id == indexInt {
	// 		model.Users = append(model.Users[:index], model.Users[index+1:]...)
	// 		fmt.Fprintf(writer, "Deleted user of id: %s", indexParam)
	// 	}
	// }

	result, _ := collection.DeleteOne(ctx, model.User{Id: id})
	json.NewEncoder(writer).Encode(result)
}
