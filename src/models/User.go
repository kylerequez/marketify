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
	Age         uint
	Gender      string
	Email       string
	Password    []byte
	Authorities []string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
