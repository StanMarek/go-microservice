package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"microservice/src/database"
	"microservice/src/model"

	"net/http"
	"time"

	uv "microservice/src/validation"

	"github.com/go-playground/validator"
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

const loginPrefix = "@new_"

func AddUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

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

	result, err := database.InsertUser(ctx, &user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(result)
}
