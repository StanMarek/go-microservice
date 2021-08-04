package handler

import (
	"context"
	"microservice/src/authentication"
	"net/http"
	"time"
)

func Logout(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	authId, _, err := authentication.ExtractTokenMetadata(request)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	deleteAuth, err := authentication.DeleteAuthentication(ctx, authId)
	if err != nil || deleteAuth == 0 {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(`{"message": "Logged out"}`))
}
