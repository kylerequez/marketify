package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/marketify/src/middlewares"
	"github.com/kylerequez/marketify/src/services"
)

type UserHandler struct {
	App        *fiber.App
	Us         *services.UserService
	Middleware *middlewares.MiddlewareHandler
}

func NewUserHandler(
	app *fiber.App,
	us *services.UserService,
	mh *middlewares.MiddlewareHandler,
) *UserHandler {
	return &UserHandler{
		App:        app,
		Us:         us,
		Middleware: mh,
	}
}

func (uh *UserHandler) Init() error {
	adminView := uh.App.Group(
		"/dashboard",
		uh.Middleware.IsLoggedIn,
		uh.Middleware.IsAdmin,
	)
	adminView.Get("/users", uh.GetUsersPage)

	return nil
}

func (uh *UserHandler) GetUsersPage(c fiber.Ctx) error {
	return uh.Us.GetUsersPage(c)
}
