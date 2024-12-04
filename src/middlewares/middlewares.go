package middlewares

import (
	"github.com/kylerequez/marketify/src/db"
	"github.com/kylerequez/marketify/src/repositories"
)

type MiddlewareHandler struct {
	Ur    *repositories.UserRepository
	Store *db.PostgresStorage
}

func NewMiddlewareHandler(ur *repositories.UserRepository, store *db.PostgresStorage) *MiddlewareHandler {
	return &MiddlewareHandler{
		Ur:    ur,
		Store: store,
	}
}
