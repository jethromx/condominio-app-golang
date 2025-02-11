package auth

import (
	"strings"

	"com.mx/crud/config"
	"com.mx/crud/internal/models"
	"com.mx/crud/internal/service"
	"com.mx/crud/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Definir constantes para los mensajes
const (
	ErrParsingInput           = "Error parsing input "
	ErrInvalidRequest         = "Invalid request "
	ErrRegisteringUser        = "Error registering user "
	ErrCouldNotRegisterUser   = "Could not register user "
	MsgUserRegistered         = "User registered successfully "
	ErrLoggingIn              = "Error logging in "
	ErrInvalidEmailOrPassword = "Invalid email or password "
	MsgSuccessLogin           = "Success login"
	ErrRefreshingToken        = "Error refreshing token "
	ErrInvalidOrExpiredToken  = "Invalid or expired refresh token "
	MsgTokenRefreshed         = "Token refreshed successfully "

	ErrInvalidOrExpiredJWT = "Invalid or expired JWT "
	ErrInvalidJWTClaims    = "Invalid JWT claims "
	ErrorParsingToken      = "Missing or malformed JWT "
)

// RegisterHandler maneja el registro de nuevos usuarios
func RegisterHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(RegisterInput)

		err := utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		user := MapUserInputToModel(input)

		user, err = authService.RegisterUser(user)
		if err != nil {
			log.Debug(ErrRegisteringUser, err)
			return utils.HandleResponse(c, fiber.StatusUnprocessableEntity, ErrCouldNotRegisterUser+err.Error(), nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, MsgUserRegistered, user)
	}
}

// LoginHandler maneja el inicio de sesión de los usuarios
func LoginHandler(authService service.AuthService, tokenService service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(LoginInput)

		err := utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		user, accessToken, refreshToken, err := authService.LoginUser(input.Identity, input.Password)

		if err != nil {
			log.Debug(ErrLoggingIn, err)
			return utils.HandleResponse(c, fiber.StatusBadRequest, ErrInvalidEmailOrPassword, nil)
		}

		if user == nil {
			return utils.HandleResponse(c, fiber.StatusForbidden, ErrInvalidEmailOrPassword, nil)
		}

		err = tokenService.RevokeAllUserTokens(user.ID)
		if err != nil {
			log.Debug(ErrLoggingIn, err)
			return utils.HandleResponse(c, fiber.StatusForbidden, ErrInvalidEmailOrPassword, nil)
		}

		var tokenEntity models.Token
		tokenEntity.Token = accessToken
		tokenEntity.RefreshToken = refreshToken
		tokenEntity.UserID = user.ID
		tokenEntity.Expirated = false
		tokenEntity.Revoked = false

		err = tokenService.SaveToken(&tokenEntity)

		if err != nil {

			log.Debug(ErrLoggingIn, err)
			return utils.HandleResponse(c, fiber.StatusForbidden, ErrInvalidEmailOrPassword, nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, MsgSuccessLogin, fiber.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

// RefreshTokenHandler maneja la renovación del token de acceso
func RefreshTokenHandler(authService service.AuthService, tokenService service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		input := new(RefreshTokenInput)

		err := utils.HandlerValidation(c, input)
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		// validar token en base de datos
		tokenEntity, err := tokenService.IsTokenValid("refreshToken", input.RefreshToken)

		if err != nil {
			log.Debug("Error validando token: ", err)
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredToken, nil)
		}

		tokenEntity.Expirated = true
		tokenEntity.Revoked = true
		err = tokenService.UpdateToken(tokenEntity)

		if err != nil {
			log.Debug("Error actualizando token: ", err)
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "error", nil)
		}

		accessToken, refreshToken, err := authService.RefreshAccessToken(input.RefreshToken)

		if err != nil {
			log.Debug("Error generando tokens :", err)
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredToken, nil)
		}

		return utils.HandleResponse(c, fiber.StatusOK, MsgTokenRefreshed, fiber.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

// LogoutHandler maneja la revocación del token de acceso
func LogoutHandler(tokenService service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		// Verificar si el encabezado de autorización está presente
		if authHeader == "" {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredToken, nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// recuperar el token JWT
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(config.GetServerSettings().JWT.SecretKey), nil
		})

		// Si hay un error al analizar el token, devolver un error
		if err != nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredJWT, nil)
		}

		// validar token en base de datos
		tokenEntity, err := tokenService.IsTokenValid("accessToken", tokenString)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredJWT, nil)
		}

		if tokenEntity == nil {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredJWT, nil)
		}

		if tokenEntity.Revoked || tokenEntity.Expirated {
			return utils.HandleResponse(c, fiber.StatusUnauthorized, ErrInvalidOrExpiredJWT, nil)
		}

		tokenEntity.Expirated = true
		tokenEntity.Revoked = true
		err = tokenService.UpdateToken(tokenEntity)

		if err != nil {
			return utils.HandleResponse(c, fiber.StatusInternalServerError, "error", nil)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
