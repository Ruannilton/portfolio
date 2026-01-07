package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         int
	AppEnv       string
	IsProduction bool

	// Database
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string
	DBSchema   string

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string

	// GitHub OAuth
	GithubClientID     string
	GithubClientSecret string

	// Session
	SessionKey    string
	SessionMaxAge int

	// JWT
	JWTSecretKey string
	JWTIssuer    string
}

func LoadConfig() (*Config, error) {
	// Carrega o arquivo .env (ignora erro se não existir em produção)
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getEnvAsInt("PORT", 8080),
		AppEnv:             getEnv("APP_ENV", "local"),
		IsProduction:       getEnvAsBool("IS_PRODUCTION", false),
		DBHost:             getEnv("BLUEPRINT_DB_HOST", ""),
		DBPort:             getEnv("BLUEPRINT_DB_PORT", "5432"),
		DBDatabase:         getEnv("BLUEPRINT_DB_DATABASE", ""),
		DBUsername:         getEnv("BLUEPRINT_DB_USERNAME", ""),
		DBPassword:         getEnv("BLUEPRINT_DB_PASSWORD", ""),
		DBSchema:           getEnv("BLUEPRINT_DB_SCHEMA", "public"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		SessionKey:         getEnv("SESSION_KEY", ""),
		SessionMaxAge:      getEnvAsInt("SESSION_MAX_AGE", 86400),
		JWTSecretKey:       getEnv("JWT_SECRET_KEY", ""),
		JWTIssuer:          getEnv("JWT_ISSUER", ""),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var errs []error

	// Validações obrigatórias do banco de dados
	if c.DBHost == "" {
		errs = append(errs, errors.New("BLUEPRINT_DB_HOST is required"))
	}
	if c.DBDatabase == "" {
		errs = append(errs, errors.New("BLUEPRINT_DB_DATABASE is required"))
	}
	if c.DBUsername == "" {
		errs = append(errs, errors.New("BLUEPRINT_DB_USERNAME is required"))
	}
	if c.DBPassword == "" {
		errs = append(errs, errors.New("BLUEPRINT_DB_PASSWORD is required"))
	}

	if c.SessionKey == "" {
		errs = append(errs, errors.New("SESSION_KEY is required"))
	}

	if c.JWTSecretKey == "" {
		errs = append(errs, errors.New("JWT_SECRET_KEY is required"))
	}

	if c.JWTIssuer == "" {
		errs = append(errs, errors.New("JWT_ISSUER is required"))
	}

	// Validações de OAuth (obrigatórias em produção)
	if c.IsProduction {
		if c.GoogleClientID == "" {
			errs = append(errs, errors.New("GOOGLE_CLIENT_ID is required in production"))
		}
		if c.GoogleClientSecret == "" {
			errs = append(errs, errors.New("GOOGLE_CLIENT_SECRET is required in production"))
		}
		if c.SessionKey == "" || c.SessionKey == "your_session_key" {
			errs = append(errs, errors.New("SESSION_KEY must be set to a secure value in production"))
		}
		if c.GithubClientID == "" {
			errs = append(errs, errors.New("GITHUB_CLIENT_ID is required in production"))
		}
		if c.GithubClientSecret == "" {	
			errs = append(errs, errors.New("GITHUB_CLIENT_SECRET is required in production"))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// getEnv retorna o valor da variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retorna o valor da variável de ambiente como int ou um valor padrão
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool retorna o valor da variável de ambiente como bool ou um valor padrão
func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
