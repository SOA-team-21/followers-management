package model

import (
	"encoding/json"
	"errors"
	"io"
)

type Person struct {
	Id      int64 `json:"id"`
	UserId  int64 `json:"userId"`
	Name    string
	Surname string
	Picture string
	Bio     string
	Quote   string
	Email   string
}

func (person *Person) Validate() error {
	if person.Name == "" {
		return errors.New("invalid name")
	}
	if person.Surname == "" {
		return errors.New("invalid surname")
	}
	if person.Email == "" {
		return errors.New("invalid email")
	}
	return nil
}

type People []*Person

func (o *Person) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Person) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *People) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
