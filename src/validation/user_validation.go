package validation

import (
	"context"
	"time"
	"unicode"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PasswordValidation(fl validator.FieldLevel) bool {
	const minLength = 6
	var upperCase bool = false
	var lowerCase bool = false
	var number bool = false
	var currentLength = 0
	password := fl.Field().String()

	for _, character := range password {
		if unicode.IsNumber(character) {
			number = true
			currentLength++
		}
		if unicode.IsUpper(character) {
			upperCase = true
			currentLength++
		}
		if unicode.IsLower(character) {
			lowerCase = true
			currentLength++
		}
	}
	if upperCase && lowerCase && number && currentLength >= minLength {
		return true
	} else {
		return false
	}
}

func UniqueEmailValidation(fl validator.FieldLevel) bool {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	client, _ := mongo.Connect(ctx, clientOptions)
	collection := client.Database("go_microservice").Collection("users")

	email := fl.Field().String()

	var foundUserByEmailBson bson.M
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUserByEmailBson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true
		}
	}
	return false
}

func UniqueLoginValidation(fl validator.FieldLevel) bool {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	client, _ := mongo.Connect(ctx, clientOptions)
	collection := client.Database("go_microservice").Collection("users")

	login := fl.Field().String()

	var user bson.M
	err := collection.FindOne(ctx, bson.M{"login": login}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true
		}
	}
	return false
}
