package api_controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CheckStatus(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"info":    "route working",
		"message": "welcome to /api/v1",
	})
}
