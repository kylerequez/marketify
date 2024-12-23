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
		"/dashboard/users",
		uh.Middleware.IsLoggedIn,
		uh.Middleware.IsAdmin,
	)
	adminView.Get("/", uh.GetUsersPage)
	adminView.Get("/:id", uh.GetUserPage)

	adminApi := uh.App.Group(
		"/api/v1/users",
		uh.Middleware.IsLoggedIn,
		uh.Middleware.IsAdmin,
	)
	adminApi.Get("/:id/edit", uh.GetUserEditForm)
	adminApi.Get("/:id/delete", uh.GetUserDeleteForm)
	adminApi.Post("/:id", uh.EditUser)
	adminApi.Delete("/:id", uh.DeleteUser)

	return nil
}

func (uh *UserHandler) GetUsersPage(c fiber.Ctx) error {
	return uh.Us.GetUsersPage(c)
}

func (uh *UserHandler) GetUserPage(c fiber.Ctx) error {
	return uh.Us.GetUserPage(c)
}

func (uh *UserHandler) GetUserEditForm(c fiber.Ctx) error {
	return uh.Us.GetUserEditForm(c)
}

func (uh *UserHandler) GetUserDeleteForm(c fiber.Ctx) error {
	return uh.Us.GetUserDeleteForm(c)
}

func (uh *UserHandler) EditUser(c fiber.Ctx) error {
	return uh.Us.EditUser(c)
}

func (uh *UserHandler) DeleteUser(c fiber.Ctx) error {
	return uh.Us.DeleteUser(c)
}
