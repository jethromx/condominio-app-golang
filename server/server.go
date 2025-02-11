package server

import (
	"context"

	"com.mx/crud/api/v0/info"

	"com.mx/crud/config"
	_ "com.mx/crud/docs" // Importa los documentos generados por swag
	"com.mx/crud/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

var GlobalServer *Server

type Server struct {
	Settings *config.ServerSettings
}

func NewServer(ctx context.Context) {
	// Obtenemos la configuración del servidor
	settings := config.GetServerSettings()

	log.Debug("Initializing fiber app")

	app := fiber.New()

	// Configurar middleware
	app.Use(cors.New())
	app.Use(requestid.New())
	app.Use(helmet.New())
	//app.Use(logger.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format:     "${locals:requestid} ${status} - ${method} ${path} \n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/Mexico_City",
	}))

	// Configurar rutas
	router.SetupRoutes(app)

	// Ruta para la documentación de Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Configurar rutas de información
	app.Get("/_service/v0/info", info.HttpInfo)

	// health check
	app.Use(healthcheck.New())

	// Iniciar el servidor
	if err := app.Listen(settings.HttpPort); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	log.Debug("Initializing app listen in port" + settings.HttpPort)

	// Variable global de acceso al servidor
	GlobalServer = &Server{
		Settings: settings,
	}
}
