package model

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty"`
	Login          string             `json:"login,omitempty" bson:"login,omitempty"`
	HashedPassword []byte             `json:"hashed_password,omitempty" bson:"hashed_password,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Activities     Activity           `json:"activity,omitempty" bson:"activity,omitempty"`
}

func (u *User) ParseEmailToLogin() string {
	var login []rune
	for _, char := range u.Email {
		if char == '@' {
			break
		}
		login = append(login, char)
	}
	return string(login)
}

func (u *User) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}
