package handler

import (
	"context"
	"encoding/json"
	"microservice/model"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)
	collection := client.Database("go_microservice").Collection("users")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	var users []model.User
	for cursor.Next(ctx) {
		var user model.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(writer).Encode(users)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)
	collection := client.Database("go_microservice").Collection("users")
	//cursor, err := collection.Find(ctx, bson.M{})
	// if err != nil {
	// 	writer.WriteHeader(http.StatusInternalServerError)
	// 	writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
	// 	return
	// }
	params := mux.Vars(request)
	//indexParam := params["id"]
	id, _ := primitive.ObjectIDFromHex(params["id"])
	// indexInt, _ := strconv.Atoi(indexParam)
	// var jsonUser model.User
	// for _, user := range model.Users {
	// 	if user.Id == indexInt {
	// 		jsonUser = user
	// 		break
	// 	}
	// }
	//json.NewEncoder(writer).Encode(jsonUser)
	var user model.User
	err := collection.FindOne(ctx, model.User{Id: id}).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	user.ToJson(writer)
}
