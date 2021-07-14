package model

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	Id            primitive.ObjectID `json:"-" bson:"_id"`
	Key           string             `json:"key" bson:"key"`
	Activity      string             `json:"activity" bson:"activity"`
	Type          string             `json:"type" bson:"type"`
	Participants  int                `json:"participants" bson:"participants"`
	Price         float64            `json:"price" bson:"price"`
	Link          string             `json:"link,omitempty" bson:"link,omitempty"`
	Accessability float64            `json:"accessability" bson:"accessability"`
}

func (a *Activity) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Activity) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
