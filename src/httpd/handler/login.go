package handler

import (
	"context"
	"encoding/json"
	auth "microservice/src/authentication"
	"microservice/src/database"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	var credentials Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	user, err := database.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(credentials.Password))
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	token, err := auth.GenerateJWT(user.Id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = auth.GenerateAuthentication(ctx, user.Id, token)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(`{"message": "Logged in"}`))
	json.NewEncoder(writer).Encode(tokens)
}
