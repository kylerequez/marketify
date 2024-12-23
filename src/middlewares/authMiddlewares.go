package middlewares

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/kylerequez/marketify/src/models"
	"github.com/kylerequez/marketify/src/shared"
)

func (mh *MiddlewareHandler) IsLoggedIn(c fiber.Ctx) error {
	log.Println("IS LOGGED IN MIDDLEWARE")

	cookie := c.Cookies("marketify-user-session", "")
	if cookie == "" {
		return c.Redirect().To("/login")
	}

	sessionId, err := uuid.Parse(cookie)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "PARSING ERROR",
			"error":   err.Error(),
		})
	}

	session, err := mh.Store.GetSession(sessionId)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "SESSION",
			"error":   err.Error(),
		})
	}

	stringUserId := string(session.Value)
	userId, err := uuid.Parse(stringUserId)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "USER ID",
			"error":   err.Error(),
		})
	}

	user, err := mh.Ur.GetUserById(c.Context(), userId)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "GET USER ID",
			"error":   err.Error(),
		})
	}

	c.Locals("loggedInUser", user)
	return c.Next()
}

func (mh *MiddlewareHandler) IsAdmin(c fiber.Ctx) error {
	log.Println("IS ADMIN MIDDLEWARE")

	user, ok := c.Locals("loggedInUser").(*models.User)
	if !ok {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "user is not logged in",
		})
	}

	isAdmin, err := mh.Ur.IsRole(c.Context(), user.ID, shared.ROLES["ADMIN"])
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
