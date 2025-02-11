package config

import (
	"os"
	"strconv"
	"sync"
)

type DatabaseSettings struct {
	Host         string `json:"DB_HOST" default:"localhost"`
	Port         int    `json:"DB_PORT" default:"5432"`
	User         string `json:"DB_USER" default:"user1"`
	Password     string `json:"DB_PASSWORD" default:"password"`
	DatabaseName string `json:"DB_DATABASE_NAME" default:"crud"`
	SslMode      string `json:"DB_SSLMODE" default:"disable"`
	TimeZone     string `json:"DB_TIMEZONE" default:"America/Mexico_City"`
}

type JWTSettings struct {
	SecretKey string `json:"JWT_SECRET" default:"DFSGHHGDHFSGGS34532"`
}

type ServerSettings struct {
	Env         string `json:"ENVIRONMENT" default:"LOCAL"`
	RuntimeMode string `json:"RUNTIME_MODE" default:"online"`
	ServiceName string `json:"SERVICE_NAME" default:"backendcore-api"`
	AppName     string `json:"APP_NAME" default:"backendcore"`
	LogLevel    string `json:"LOG_LEVEL" default:"DEBUG"`
	HttpPort    string `json:"HTTP_PORT" default:"8080"`

	Database DatabaseSettings

	JWT JWTSettings
}

var (
	instance *ServerSettings
	once     sync.Once
)

func GetServerSettings() *ServerSettings {
	once.Do(func() {

		port := os.Getenv("SERVER_PORT")
		if port == "" {
			port = "8080"
		}

		runtimeMode := os.Getenv("RUNTIME_MODE")
		if runtimeMode == "" {
			runtimeMode = "online"
		}

		instance = &ServerSettings{
			Env:         os.Getenv("ENV"),
			ServiceName: os.Getenv("SERVICE_NAME"),
			AppName:     os.Getenv("APP_NAME"),
			LogLevel:    os.Getenv("LOG_LEVEL"),

			HttpPort:    ":" + port,
			RuntimeMode: runtimeMode,

			Database: DatabaseSettings{
				Host:         getEnv("DB_HOST", "localhost"),
				Port:         getEnvAsInt("DB_PORT", 5432),
				User:         getEnv("DB_USER", "user1"),
				Password:     getEnv("DB_PASSWORD", "password"),
				DatabaseName: getEnv("DB_DATABASE", "crud"),
				SslMode:      getEnv("DB_SSLMODE", "disable"),
				TimeZone:     getEnv("DB_TIMEZONE", "America/Mexico_City"),
			},
			JWT: JWTSettings{
				SecretKey: getEnv("JWT_SECRET", "DFSGHHGDHFSGGS34532"),
			},
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
