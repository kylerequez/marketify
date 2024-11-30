package utils

import "github.com/gofiber/fiber/v3"

func Redirect(c fiber.Ctx, status int, url string) error {
	c.Set("HX-Redirect", url)
	return c.SendStatus(status)
}
