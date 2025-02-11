package main

import (
	"context"

	"com.mx/crud/config/database"
	"com.mx/crud/server"
	"github.com/gofiber/fiber/v2/log"
)

// @title Condominium Management API
// @version 1.0
// @description API for managing condominiums, buildings, apartments, and residents.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {

	log.Info("Initializing the application")
	// Generamos un contexto
	ctx := context.Background()

	// connect to the database
	database.ConnectDatabase()
	// close the database connection when the main function ends
	defer database.Close()

	// Generamos el servidor
	server.NewServer(ctx)

}
