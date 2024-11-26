package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/marketify/src/services"
)

type UserHandler struct {
	App *fiber.App
	Us  *services.UserService
}

func NewUserHandler(app *fiber.App, us *services.UserService) *UserHandler {
	return &UserHandler{
		App: app,
		Us:  us,
	}
}

func (uh *UserHandler) Init() error {
	return nil
}
