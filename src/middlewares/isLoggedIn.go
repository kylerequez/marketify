package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func IsLoggedIn(c fiber.Ctx) error {
	cookie := c.Cookies("marketify-user-session")

	log.Println("COOKIE: ", cookie)

	return c.Next()
}
