package database

import (
	"microservice/src/model"

	"go.mongodb.org/mongo-driver/mongo"
)

var activityCollection *mongo.Collection

func InsertActivity(activity *model.Activity) (*mongo.InsertOneResult, error) {
	result, err := activityCollection.InsertOne(ctx, activity)
	return result, err
}
