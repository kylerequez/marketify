package middlewares

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/kylerequez/marketify/src/shared"
)

func (mh *MiddlewareHandler) IsLoggedIn(c fiber.Ctx) error {
	log.Println("IS LOGGED IN MIDDLEWARE")

	cookie := c.Cookies("marketify-user-session", "")
	if cookie == "" {
		return c.Redirect().To("/login")
	}

	c.Locals("userId", cookie)
	return c.Next()
}

func (mh *MiddlewareHandler) IsAdmin(c fiber.Ctx) error {
	log.Println("IS LOGGED IN MIDDLEWARE")

	id := c.Locals("userId")

	userId, err := uuid.Parse(id.(string))
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isAdmin, err := mh.Ur.IsRole(c.Context(), userId, shared.ROLES["ADMIN"])
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !isAdmin {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "it is not an admin",
		})
	}

	return c.Next()
}
