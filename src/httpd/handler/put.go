package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"microservice/src/database"

	"microservice/src/model"
	uv "microservice/src/validation"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUpdateRequest struct {
	NewEmail    string `json:"email" validate:"email,uni_email"`
	NewLogin    string `json:"login" validate:"login"`
	NewPassword string `json:"password" validate:"password"`
}

func (uur *UserUpdateRequest) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("password", uv.PasswordValidation)
	validate.RegisterValidation("uni_email", uv.UniqueEmailValidation)
	validate.RegisterValidation("login", uv.UniqueLoginValidation)
	return validate.Struct(uur)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	var updatedUser UserUpdateRequest
	json.NewDecoder(request.Body).Decode(&updatedUser)
	err := updatedUser.Validate()
	if err != nil {
		http.Error(
			writer,
			fmt.Sprintf("Error validating user udpdate request: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user model.User

	user.Email = updatedUser.NewEmail
	if strings.Contains(updatedUser.NewLogin, loginPrefix) {
		fmt.Fprint(writer, "Updated login contains prefix for new default logins. Please ensure that your new login does not contain that prefix")
		return
	}
	if updatedUser.NewLogin != "" {
		user.Login = updatedUser.NewLogin
	} else {
		user.Login = user.ParseEmailToLogin()
	}
	user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(updatedUser.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user.UpdatedAt = time.Now()

	result, err := database.UpdateUserById(ctx, id, user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(result)
}
