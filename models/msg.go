package models

import (
	"fmt"
	"net/http"
)

type Msg struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Body string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type MsgList struct {
	Msgs []Msg `json:"msgs"`
}

func (i *Msg) Bind(r *http.Request) error {
	if i.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}

func (*MsgList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Msg) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

