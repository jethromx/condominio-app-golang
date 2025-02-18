package auth

/*
import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
)

// Redirigir al login de Google
func GoogleLogin(c *fiber.Ctx) error {
	url := auth.GetGoogleAuthURL()
	return c.Redirect(url)
}

// Callback de Google
func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := auth.GetGoogleUserInfo(code)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No se pudo obtener el token"})
	}

	// Obtener info del usuario
	client := google.Config{}.Client(c.Context(), token)
	service, err := oauth2.New(client)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error en el servicio OAuth"})
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo obtener info del usuario"})
	}

	// Verificar si el usuario ya existe en la BD
	var usuario models.Usuario
	if err := config.DB.Where("email = ?", userInfo.Email).First(&usuario).Error; err != nil {
		usuario = models.Usuario{
			Email:  userInfo.Email,
			Nombre: userInfo.Name,
			Rol:    "residente",
		}
		config.DB.Create(&usuario)
	}

	// Generar JWT
	jwtToken, err := auth.GenerateJWT(usuario.ID, usuario.Email, usuario.Rol)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error generando token"})
	}

	return c.JSON(fiber.Map{"token": jwtToken})
}
*/
