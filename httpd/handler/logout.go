package handler

import (
	"microservice/authentication"
	"net/http"
)

func Logout(writer http.ResponseWriter, request *http.Request) {
	authId, _, err := authentication.ExtractTokenMetadata(request)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	deleteAuth, err := authentication.DeleteAuthentication(authId)
	if err != nil || deleteAuth == 0 {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"message": "Logged out"}`))
}
