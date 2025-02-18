package utils

import "github.com/gofiber/fiber/v2"

// handleResponse es una función auxiliar para manejar las respuestas y los errores
func HandleResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	})
}

func DataPagination(page int, pageSize int, totalRecords int, data interface{}) fiber.Map {
	return fiber.Map{
		"page":         page,
		"pageSize":     pageSize,
		"totalPages":   (int64(totalRecords) + int64(pageSize) - 1) / int64(pageSize),
		"totalRecords": totalRecords,
		"items":        data,
	}
}
