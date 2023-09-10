package pages_controllers

import "github.com/gofiber/fiber/v2"

func HomePage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}
