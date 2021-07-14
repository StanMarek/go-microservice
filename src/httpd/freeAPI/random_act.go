package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"microservice/src/database"
	"microservice/src/model"
	"net/http"

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

	response, err := http.Get("https://www.boredapi.com/api/activity")
	if err != nil {
		log.Fatal(err)
	}

	reponseData, _ := ioutil.ReadAll(response.Body)
	var responseObj ResponseData
	json.Unmarshal(reponseData, &responseObj)
	activity := responseObj.ParseActivity()
	activity.Id = primitive.NewObjectID()
	result, err := database.InsertActivity(activity)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": ` + err.Error() + `"}`))
		return

	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(result)
	json.NewEncoder(res).Encode(activity)
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
