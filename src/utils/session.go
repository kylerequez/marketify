package utils

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/kylerequez/marketify/src/models"
	"github.com/kylerequez/marketify/src/repositories"
)

func RetrieveLoggedInUser(c fiber.Ctx, ur *repositories.UserRepository) *models.User {
	cookie := c.Cookies("marketify-user-session")

	id, err := uuid.Parse(cookie)
	if err != nil {
		log.Println("err: ", err)
		return nil
	}

	user, err := ur.GetUserById(c.Context(), id)
	if err != nil {
		log.Println("err: ", err)
		return nil
	}

	return user
}
