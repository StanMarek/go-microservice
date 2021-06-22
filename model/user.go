package model

import (
	"encoding/json"
	"io"
	"time"
)

// TODO: add more validation rules
type User struct {
	Id        int       `json:"_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Login     string    `json:"login,omitempty"`
	Password  string    `json:"password,ozmitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

var Users []User

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

func NextId() int {
	return len(Users) + 1
}

func (u *User) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(u)
}

func Exists(id int) (bool, int) {
	for index, i := range Users {
		if id == i.Id {
			return true, index
		}
	}
	return false, -1
}
