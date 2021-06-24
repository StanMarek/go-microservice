package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"microservice/model"
	"net/http"
	"time"

	uv "microservice/validation"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserPostRequest struct {
	Email string `json:"email" validate:"email,required"`
	//Login 		string	`json:"login" validate:"required"`
	Password string `json:"password" validate:"password,required"`
}

func (upr *UserPostRequest) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", uv.PasswordValidation)
	return validate.Struct(upr)
}

var client *mongo.Client

func AddUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)

	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// client, err := mongo.Connect(context.TODO(), clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = client.Ping(context.TODO(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Connected to mongo")

	// writer.Header().Add("content-type", "application/json")

	var userPostRequest UserPostRequest
	json.NewDecoder(request.Body).Decode(&userPostRequest)
	err := userPostRequest.Validate()
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Error validating user: %s", err),
			http.StatusBadRequest,
		)
		return
	}

	var user model.User
	user.CreatedAt = time.Now()
	// user.Id = model.NextId()
	user.Email = userPostRequest.Email
	user.Login = user.ParseEmailToLogin()
	user.Password = userPostRequest.Password

	// exists, _ := model.Exists(user.Id)
	// if exists {
	// 	fmt.Fprintf(writer, "User of id {%d} already exists", user.Id)
	// 	return
	// } else {
	// 	model.Users = append(model.Users, user)
	// 	//json.NewEncoder(writer).Encode(user)
	// 	user.ToJson(writer)
	// }

	collection := client.Database("go_microservice").Collection("users")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(writer).Encode(result)
}
