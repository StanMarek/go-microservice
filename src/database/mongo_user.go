package database

import (
	"context"
	"microservice/src/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var clientOptions *options.ClientOptions
var collection *mongo.Collection

func Connect() error {
	clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	collection = client.Database("go_microservice").Collection("users")
	activityCollection = client.Database("go_microservice").Collection("activities")
	return nil
}

func InsertUser(ctx context.Context, user *model.User) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, user)
	return result, err
}

func GetUsers(ctx context.Context) ([]model.User, error) {
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

func GetUserById(ctx context.Context, id primitive.ObjectID) (model.User, error) {
	var user model.User

	err := collection.FindOne(ctx, model.User{Id: id}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	err := collection.FindOne(ctx, model.User{Login: login}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func UpdateUserById(ctx context.Context, id primitive.ObjectID, user model.User) (*mongo.UpdateResult, error) {
	update := bson.D{
		{"$set", bson.D{
			{"email", user.Email},
			{"login", user.Login},
			{"hashed_password", user.HashedPassword},
			{"updated_at", user.UpdatedAt},
		}}}
	result, err := collection.UpdateByID(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteUserById(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(ctx, model.User{Id: id})
	if err != nil {
		return nil, err
	}
	return result, nil
}
