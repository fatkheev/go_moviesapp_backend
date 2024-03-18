package model

import (
	"filmoteca/services"
)

type Actor struct {
	ID        int               `json:"id,omitempty"`
	Name      string            `json:"name"`
	Gender    string            `json:"gender,omitempty"`
	Birthdate services.JSONDate `json:"birthdate"`
	Movies    []Movie           `json:"movies,omitempty"`
}
