package middleware

import "github.com/gofiber/fiber/v2"

// Middleware para verificar el rol del usuario
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRol := c.Locals("userRol")

		if userRol != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Acceso denegado"})
		}

		return c.Next()
	}
}
