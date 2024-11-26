package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/marketify/src/services"
)

type authHandler interface {
	Init() error
}

type AuthHandler struct {
	App *fiber.App
	Us  *services.UserService
}

func NewAuthHandler(app *fiber.App, us *services.UserService) *AuthHandler {
	return &AuthHandler{
		App: app,
		Us:  us,
	}
}

func (ah *AuthHandler) Init() error {
	api := ah.App.Group("/api/v1/auth")
	api.Post("/login", ah.LoginUser)

	return nil
}

func (ah *AuthHandler) LoginUser(c fiber.Ctx) error {
	return ah.Us.LoginUser(c)
}
