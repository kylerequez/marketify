package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Firstname   string
	Middlename  string
	Lastname    string
	Birthdate   time.Time
	Gender      string
	Email       string
	Password    []byte
	Authorities []string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
