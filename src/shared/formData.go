package shared

import (
	"github.com/google/uuid"
)

type LoginFormData struct {
	Email    string
	Password string
	Errors   map[string]string
}

type SignupFormData struct {
	Firstname  string
	Middlename string
	Lastname   string
	Birthdate  string
	Gender     string
	Email      string
	Password   string
	RePassword string
	Errors     map[string]string
}

type EditUserFormData struct {
	ID         uuid.UUID
	Firstname  string
	Middlename string
	Lastname   string
	Gender     string
	Email      string
	Errors     map[string]string
}
