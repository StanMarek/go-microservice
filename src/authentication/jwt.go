package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"microservice/src/database"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenMetadata struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessId     string `json:"access_id"`
	RefreshId    string `json:"refresh_id"`
	AccessExp    int64  `json:"access_exp"`
	RefreshExp   int64  `json:"refresh_exp"`
}

// https://pkg.go.dev/github.com/golang-jwt/jwt#section-documentation
func GenerateJWT(userID primitive.ObjectID) (*TokenMetadata, error) {
	// $ set ACCESS_TOKEN=jwtsecrettoken
	// os.Setenv("ACCESS_TOKEN", "jwtsecrettoken")
	tokenData := &TokenMetadata{}
	tokenData.AccessExp = time.Now().Add(time.Minute * 5).Unix()
	tokenData.AccessId = uuid.NewV4().String()

	tokenData.RefreshExp = time.Now().Add(time.Hour * 24).Unix()
	tokenData.AccessId = uuid.NewV4().String()

	// var signingAccessKey = os.Getenv("ACCESS_TOKEN")
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["user_id"] = userID
	accessTokenClaims["access_id"] = tokenData.AccessId
	accessTokenClaims["access_exp"] = tokenData.AccessExp

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	var err error
	tokenData.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN")))
	if err != nil {
		return nil, err
	}

	// set REFRESH_TOKEN=jwtrefreshtoken
	// os.Setenv("REFRESH_TOKEN", "jwtrefreshtoken")
	// var signingRefreshKey = os.Getenv("REFRESH_TOKEN")
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["user_id"] = userID
	refreshTokenClaims["refresh_id"] = tokenData.RefreshId
	refreshTokenClaims["refresh_exp"] = tokenData.RefreshId

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenData.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN")))
	if err != nil {
		return nil, err
	}
	return tokenData, nil
}

func ExtractToken(request *http.Request) string {
	token := request.Header.Get("Authorization")
	arr := strings.Split(token, " ")
	if len(arr) == 2 {
		return arr[1]
	}
	return ""
}

func Verify(request *http.Request) (*jwt.Token, error) {
	tokenStr := ExtractToken(request)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, isValid := token.Method.(*jwt.SigningMethodHMAC)
		if !isValid {
			return nil, fmt.Errorf("%v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func IsTokenValid(request *http.Request) error {
	token, err := Verify(request)
	if err != nil {
		return err
	}
	_, isValid := token.Claims.(jwt.Claims)
	if !isValid && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(request *http.Request) (string, primitive.ObjectID, error) {
	token, err := Verify(request)
	if err != nil {
		return "", primitive.NilObjectID, err
	}
	claims, isValid := token.Claims.(jwt.MapClaims)
	if isValid && token.Valid {
		accessId, isValid := claims["access_id"].(string)
		if !isValid {
			return "", primitive.NilObjectID, err
		}
		userIdClaims := claims["user_id"]
		userId := userIdClaims.(primitive.ObjectID)
		return accessId, userId, nil
	}
	return "", primitive.NilObjectID, err
}

func Fetch(ctx context.Context, accessId string) (primitive.ObjectID, error) {
	userIdRedis, err := database.ClientRedis.Get(ctx, accessId).Result()
	if err != nil {
		return primitive.NilObjectID, err
	}
	userId, err := primitive.ObjectIDFromHex(userIdRedis)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return userId, nil
}

func RefreshToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	ctx := request.Context()

	mapToken := map[string]string{}
	err := json.NewDecoder(request.Body).Decode(&mapToken)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	refreshToken := mapToken["refresh_id"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("bad signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN")), nil
	})

	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	if _, isValid := token.Claims.(jwt.MapClaims); !isValid && !token.Valid {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	claims, isValid := token.Claims.(jwt.MapClaims)
	if isValid && token.Valid {
		refreshId, ok := claims["refresh_id"].(string)
		if !ok {
			writer.WriteHeader(http.StatusUnprocessableEntity)
			writer.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		userIdClaims := claims["user_id"]
		userId := userIdClaims.(primitive.ObjectID)
		delete, err := DeleteAuthentication(ctx, refreshId)
		if delete == 0 || err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		newToken, err := GenerateJWT(userId)
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		err = GenerateAuthentication(ctx, userId, newToken)
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		newTokens := map[string]string{
			"access_id":  newToken.AccessId,
			"refresh_id": newToken.RefreshId,
		}
		writer.WriteHeader(http.StatusForbidden)
		json.NewEncoder(writer).Encode(newTokens)
	} else {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Expired"))
	}
}
