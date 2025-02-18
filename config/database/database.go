package database

import (
	"fmt"
	"sync"
	"time"

	"com.mx/crud/config"
	"com.mx/crud/internal/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once //La variable once de tipo sync.Once garantiza que la inicialización de la conexión a la base de datos ocurra solo una vez.
)

func ConnectDatabase() {
	once.Do(func() {
		// obtain the server settings
		settings := config.GetServerSettings()
		dsn := buildDSN(settings)

		log.Debug("Connecting to database")
		database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
			panic("failed to connect database")

		}

		// Configurar el pool de conexiones
		sqlDB, err := database.DB()
		if err != nil {
			log.Fatal("Failed to get database connection:", err)
		}

		sqlDB.SetMaxIdleConns(10)                  // Número máximo de conexiones inactivas en el pool
		sqlDB.SetMaxOpenConns(100)                 // Número máximo de conexiones abiertas
		sqlDB.SetConnMaxLifetime(time.Hour)        // Tiempo máximo de vida de una conexión
		sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Tiempo máximo de inactividad de una conexión

		// Realizar migraciones automáticas
		err = database.AutoMigrate(
			&models.Condominium{},
			&models.Building{},
			&models.Apartment{},
			&models.Resident{},
			&models.Payment{},
			&models.Maintenance{},
			&models.Reservation{},
			&models.User{},
			&models.Token{},
		)
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}

		DB = database
	})
}

// buildDSN construye el Data Source Name (DSN) para la conexión a la base de datos
func buildDSN(settings *config.ServerSettings) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=%s",
		settings.Database.Host,
		settings.Database.User,
		settings.Database.Password,
		settings.Database.DatabaseName,
		settings.Database.Port,
		settings.Database.TimeZone)
}

func Close() {
	db, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}
