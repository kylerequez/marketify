package shared

import (
	"time"
)

type UserSession struct {
	ID         string
	Value      []byte
	Expiration time.Time
}

const USER_SESSION_LIFESPAN = (15 * time.Minute)
