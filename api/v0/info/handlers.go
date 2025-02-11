package info

import (
	"github.com/gofiber/fiber/v2"
)

func HttpInfo(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	c.Status(fiber.StatusOK)
	return c.JSON(getAppInfo())
}
