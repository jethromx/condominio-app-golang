package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var (
	// AuditUserIDKey es la clave para el auditUserID en el contexto
	auditUserIDKey           = "email"
	errorAuditUserIDNotFound = errors.New("Audit user ID not found in context")
)

// GetAuditUserID extrae el auditUserID del contexto
func GetAuditUserID(c *fiber.Ctx) (string, error) {
	auditUserID, ok := c.Locals(auditUserIDKey).(string)
	if !ok {
		log.Debug(errorAuditUserIDNotFound)
		return "", errorAuditUserIDNotFound
	}
	return auditUserID, nil
}
