package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Provider struct {
	App      *App
	Database *Database
	Auth     *Auth
}

type App struct {
	Env      string
	Port     string
	Name     string
	Timezone string
}

type Database struct {
	URL string
}

type Auth struct {
	Mode                string
	FirebaseCredentials string
	JWTSecret           string
	JWTLifetimeMinutes  int
}

func Load() *Provider {
	_ = godotenv.Load()

	jwtLifetime, _ := strconv.Atoi(getEnv("JWT_LIFETIME_MINUTES", "60"))

	return &Provider{
		App: &App{
			Env:      getEnv("APP_ENV", "development"),
			Port:     getEnv("PORT", "8080"),
			Name:     getEnv("APP_NAME", "go-boilerplate"),
			Timezone: getEnv("TIMEZONE", "UTC"),
		},
		Database: &Database{
			URL: getEnv("DATABASE_URL", ""),
		},
		Auth: &Auth{
			Mode:                getEnv("AUTH_MODE", "firebase"),
			FirebaseCredentials: getEnv("FIREBASE_CREDENTIALS", ""),
			JWTSecret:           getEnv("JWT_SECRET", ""),
			JWTLifetimeMinutes:  jwtLifetime,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
