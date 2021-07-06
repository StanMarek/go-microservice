package handler

import (
	"context"
	"encoding/json"
	"log"
	"microservice/database"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(writer http.ResponseWriter, request *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	user, err := database.GetUserByLogin(credentials.Login)
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

	token, err := GenerateJWT(user.Id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = GenerateAuthentication(user.Id, token)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}

	tokens := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	json.NewEncoder(writer).Encode(tokens)
}

type TokenMetadata struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessId     string `json:"access_id"`
	RefreshId    string `json:"refresh_id"`
	AccessExp    int64  `json:"access_exp"`
	RefreshExp   int64  `json:"refresh_exp"`
}

//https://pkg.go.dev/github.com/golang-jwt/jwt#section-documentation
func GenerateJWT(userID primitive.ObjectID) (*TokenMetadata, error) {
	//$ set ACCESS_TOKEN=jwtsecrettoken
	//os.Setenv("ACCESS_TOKEN", "jwtsecrettoken")
	tokenData := &TokenMetadata{}
	tokenData.AccessExp = time.Now().Add(time.Minute * 5).Unix()
	tokenData.AccessId = uuid.NewV4().String()

	tokenData.RefreshExp = time.Now().Add(time.Hour * 24).Unix()
	tokenData.AccessId = uuid.NewV4().String()

	//var signingAccessKey = os.Getenv("ACCESS_TOKEN")
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["user_id"] = userID
	accessTokenClaims["access_token_id"] = tokenData.AccessId
	accessTokenClaims["expires"] = tokenData.AccessExp

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	var err error
	tokenData.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN")))
	if err != nil {
		return nil, err
	}

	// set REFRESH_TOKEN=jwtrefreshtoken
	//os.Setenv("REFRESH_TOKEN", "jwtrefreshtoken")
	//var signingRefreshKey = os.Getenv("REFRESH_TOKEN")
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["user_id"] = userID
	refreshTokenClaims["access_token_id"] = tokenData.RefreshId
	refreshTokenClaims["expires"] = tokenData.RefreshId

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenData.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN")))
	if err != nil {
		return nil, err
	}
	return tokenData, nil
}

var clientRedis *redis.Client
var ctxRedis = context.Background()

func ConnectRedis() {
	redisdsn := os.Getenv("REDIS_DSN")
	if len(redisdsn) == 0 {
		redisdsn = "127.0.0.1:6379"
	}
	clientRedis = redis.NewClient(&redis.Options{
		Addr:     redisdsn,
		Password: "",
		DB:       0,
	})
	if _, err := clientRedis.Ping(ctxRedis).Result(); err != nil {
		log.Fatal(err)
	}
}

func GenerateAuthentication(userId primitive.ObjectID, tokenData *TokenMetadata) error {

	accessToken := time.Unix(tokenData.AccessExp, 0)
	refreshToken := time.Unix(tokenData.RefreshExp, 0)
	now := time.Now()

	err := clientRedis.Set(ctxRedis, tokenData.AccessId, userId.String(), accessToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = clientRedis.Set(ctxRedis, tokenData.RefreshId, userId.String(), refreshToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}
