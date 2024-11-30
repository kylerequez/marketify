package middlewares

import (
	"github.com/kylerequez/marketify/src/repositories"
)

type MiddlewareHandler struct {
	Ur *repositories.UserRepository
}

func NewMiddlewareHandler(ur *repositories.UserRepository) *MiddlewareHandler {
	return &MiddlewareHandler{
		Ur: ur,
	}
}
