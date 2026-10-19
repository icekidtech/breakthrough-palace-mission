package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env            string
	Port           string
	DatabaseURL    string
	JWTSecret      string
	FrontendURL    string
	CookieDomain   string
	CookieSecure   bool
	CookieSameSite string
}

func Load() *Config {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	env := os.Getenv("ENV")
	isDev := env == "development"

	return &Config{
		Env:            env,
		Port:           os.Getenv("PORT"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		FrontendURL:    os.Getenv("FRONTEND_URL"),
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
		CookieSecure:   !isDev,
		CookieSameSite: "Lax",
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{Env: %s, Port: %s, DatabaseURL: %s, CookieSecure: %v, CookieDomain: %s}",
		c.Env, c.Port, c.DatabaseURL, c.CookieSecure, c.CookieDomain)
}
