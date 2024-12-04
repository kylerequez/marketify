package utils

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/marketify/src/models"
)

func RetrieveLoggedInUser(
	c fiber.Ctx,
) *models.User {
	user, ok := c.Locals("loggedInUser").(*models.User)
	if !ok && user == nil {
		return nil
	}

	return user
}
