package model

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	Id            primitive.ObjectID `json:"-" bson:"_id"`
	Key           string             `json:"key,omitempty" bson:"key,omitempty"`
	Activity      string             `json:"activity,omitempty" bson:"activity,omitempty"`
	Type          string             `json:"type,omitempty" bson:"type,omitempty"`
	Participants  int                `json:"participants,omitempty" bson:"participants,omitempty"`
	Price         float64            `json:"price,omitempty" bson:"price,omitempty"`
	Link          string             `json:"link,omitempty" bson:"link,omitempty"`
	Accessability float64            `json:"accessability,omitempty" bson:"accessability,omitempty"`
}

func (a *Activity) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Activity) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
