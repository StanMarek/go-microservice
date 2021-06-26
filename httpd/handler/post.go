package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"microservice/database"
	"microservice/model"
	"net/http"
	"time"

	uv "microservice/validation"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserPostRequest struct {
	Email    string `json:"email" validate:"email,required,uni_email"`
	Password string `json:"password" validate:"password,required"`
}

func (upr *UserPostRequest) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", uv.PasswordValidation)
	validate.RegisterValidation("uni_email", uv.UniqueEmailValidation)
	return validate.Struct(upr)
}

var client *mongo.Client

const loginPrefix = "@new_"

func AddUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")

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
	user.Email = userPostRequest.Email
	user.Login = loginPrefix + user.ParseEmailToLogin()
	user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(userPostRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return
	}

	result, err := database.InsertUser(&user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(writer).Encode(result)
}
