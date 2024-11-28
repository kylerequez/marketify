package shared

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID         uuid.UUID
	Value      string
	Expiration time.Time
}

const USER_SESSION_LIFESPAN = (15 * time.Minute)
