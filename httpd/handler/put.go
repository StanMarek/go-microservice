package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"microservice/model"
	uv "microservice/validation"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserUpdateRequest struct {
	Email    string `json:"email" validate:"email,required"`
	NewLogin string `json:"new_login"`
	Password string `json:"password" validate:"password,required"`
}

func (uur *UserUpdateRequest) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", uv.PasswordValidation)
	return validate.Struct(uur)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)

	// indexParam := params["id"]
	// indexInt, _ := strconv.Atoi(indexParam)

	var updatedUser UserUpdateRequest
	json.NewDecoder(request.Body).Decode(&updatedUser)
	err := updatedUser.Validate()
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Error validating user udpdate request: %s", err),
			http.StatusBadRequest,
		)
		return
	}

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user model.User
	// var id int
	// for index, userIter := range model.Users {
	// 	if userIter.Id == indexInt {
	// 		user = userIter
	// 		id = index
	// 	}
	// }

	user.Email = updatedUser.Email
	//user.Login = user.ParseEmailToLogin()
	if updatedUser.NewLogin != "" {
		user.Login = updatedUser.NewLogin
	} else {
		user.Login = user.ParseEmailToLogin()
	}
	user.Password = updatedUser.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return
	}
	user.HashedPassword = string(hashedPassword)
	//user.CreatedAt = model.Users[id].CreatedAt
	user.UpdatedAt = time.Now()

	collection := client.Database("go_microservice").Collection("users")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	// filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{{"email", user.Email},
			{"login", user.Login},
			{"password", user.Password},
			{"updated_at", user.UpdatedAt},
		},
		},
	}
	result, err := collection.UpdateByID(ctx, id, update)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	// model.Users[id] = user
	// model.Users[id].ToJson(writer)
	json.NewEncoder(writer).Encode(result)

}
