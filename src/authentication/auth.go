package authentication

import (
	"microservice/src/database"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateAuthentication(userId primitive.ObjectID, tokenData *TokenMetadata) error {

	accessToken := time.Unix(tokenData.AccessExp, 0)
	refreshToken := time.Unix(tokenData.RefreshExp, 0)
	now := time.Now()

	err := database.ClientRedis.Set(database.CtxRedis, tokenData.AccessId, userId.String(), accessToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = database.ClientRedis.Set(database.CtxRedis, tokenData.RefreshId, userId.String(), refreshToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteAuthentication(authId string) (int64, error) {
	del, err := database.ClientRedis.Del(database.CtxRedis, authId).Result()
	if err != nil {
		return 0, err
	}
	return del, nil
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		isTokenValidError := IsTokenValid(request)
		if isTokenValidError != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"message": ` + isTokenValidError.Error() + `"}`))
			return
		}
		next.ServeHTTP(writer, request)
	})
}
