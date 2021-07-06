package database

import (
	"context"
	"microservice/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var clientOptions *options.ClientOptions
var ctx context.Context
var collection *mongo.Collection

func Connect() error {
	clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	collection = client.Database("go_microservice").Collection("users")
	return nil
}

func InsertUser(user *model.User) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, user)
	return result, err
}

func GetUsers() ([]model.User, error) {
	var users []model.User

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user model.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserById(id primitive.ObjectID) (model.User, error) {
	var user model.User

	err := collection.FindOne(ctx, model.User{Id: id}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func GetUserByLogin(login string) (model.User, error) {
	var user model.User

	err := collection.FindOne(ctx, model.User{Login: login}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func UpdateUserById(id primitive.ObjectID, newData *bson.D) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateByID(ctx, id, newData)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteUserById(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(ctx, model.User{Id: id})
	if err != nil {
		return nil, err
	}
	return result, nil
}
