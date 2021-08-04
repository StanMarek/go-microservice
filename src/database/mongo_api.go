package database

import (
	"context"
	"microservice/src/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var activityCollection *mongo.Collection

func InsertActivity(ctx context.Context, activity *model.Activity) (*mongo.InsertOneResult, error) {
	var err error
	var findResult interface{}

	if err = activityCollection.FindOne(ctx, bson.M{"key": activity.Key}).Decode(&findResult); err != nil {
		result, err := activityCollection.InsertOne(ctx, activity)
		return result, err
	}

	return nil, err
}

func InsertActivityIntoUser(ctx context.Context, act model.Activity, userId primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := bson.M{"$and": []interface{}{bson.M{"_id": userId}}}
	insert := bson.M{"$set": bson.M{
		"activity.key":           act.Key,
		"activity.activity":      act.Activity,
		"activity.type":          act.Type,
		"activity.participants":  act.Participants,
		"activity.price":         act.Price,
		"activity.link":          act.Link,
		"activity.accessability": act.Accessability}}
	result, err := collection.UpdateOne(ctx, filter, insert)
	return result, err
}

func GetActivity(ctx context.Context, key string) (model.Activity, error) {
	var act model.Activity

	err := activityCollection.FindOne(ctx, bson.M{"key": key}).Decode(&act)
	if err != nil {
		return model.Activity{}, err
	}
	return act, nil
}
