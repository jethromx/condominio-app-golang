package middleware

import (
	"strings"

	"com.mx/crud/config"
	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Definir constantes para los mensajes
const (
	ErrInvalidOrExpiredJWT = "Invalid or expired JWT"
	ErrInvalidJWTClaims    = "Invalid JWT claims"
	ErrorParsingToken      = "Missing or malformed JWT"
	ErrorGettingToken      = "Error getting token"
	ErrorValidatingToken   = "Error validating token"
)

// Protected es un middleware que protege las rutas verificando el token JWT
func Protected(tokenService service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		// Verificar si el encabezado de autorización está presente
		if authHeader == "" {
			log.Debug("Missing header with token")
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrorParsingToken, nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verificar si el token está presente
		if len(tokenString) <= 0 {
			log.Debug("Error getting token")
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrorGettingToken, nil)
		}

		// recuperar el token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Debug("Unexpected signing method")
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(config.GetServerSettings().JWT.SecretKey), nil
		})

		// Si hay un error al analizar el token, devolver un error
		if err != nil {
			log.Info("Error parsing token: ", err)
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredJWT, nil)
		}

		// validar token en base de datos
		tokenEntity, err := tokenService.IsTokenValid("accessToken", tokenString)

		if err != nil {
			log.Debug("Error to validated token: ", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, ErrorValidatingToken, nil)
		}

		if tokenEntity == nil {
			log.Info("Token not found")
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredJWT, nil)
		}

		if tokenEntity.Revoked || tokenEntity.Expirated {
			log.Info("Token revoked or expired")
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredJWT, nil)
		}

		// recuperar los claims del token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user_id", claims["user_id"])
			c.Locals("username", claims["username"])
			c.Locals("email", claims["email"])
		} else {
			log.Debug("Invalid JWT claims")
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidJWTClaims, nil)
		}

		return c.Next()
	}
}
