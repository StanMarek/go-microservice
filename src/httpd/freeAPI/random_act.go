package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"microservice/src/database"
	"microservice/src/model"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseData struct {
	Key           string  `json:"key"`
	Activity      string  `json:"activity"`
	Type          string  `json:"type"`
	Participants  int     `json:"participants"`
	Price         float64 `json:"price"`
	Link          string  `json:"link"`
	Accessability float64 `json:"accessability"`
}

func RandomActvity(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	responseObj := ReadApi()
	activity := responseObj.ParseActivity()
	activity.Id = primitive.NewObjectID()
	result, err := database.InsertActivity(ctx, activity)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return

	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(result)
	activity.ToJson(res)
}
func UserRandomActivity(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	key := params["key"]

	fmt.Println("1")
	act, err := database.GetActivity(ctx, key)
	if err != nil {
		fmt.Println("3")
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	act.ToJson(res)
	fmt.Println("2")

	result, err := database.InsertActivityIntoUser(ctx, act, id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(result)
}

func ReadApi() *ResponseData {
	response, err := http.Get("https://www.boredapi.com/api/activity")
	if err != nil {
		log.Fatal(err)
	}
	reponseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObj ResponseData
	json.Unmarshal(reponseData, &responseObj)

	return &responseObj
}

func (resData *ResponseData) ParseActivity() *model.Activity {
	return &model.Activity{
		Key:           resData.Key,
		Activity:      resData.Activity,
		Type:          resData.Type,
		Participants:  resData.Participants,
		Price:         math.Round(resData.Price * 100 / 100),
		Link:          resData.Link,
		Accessability: resData.Accessability,
	}
}
