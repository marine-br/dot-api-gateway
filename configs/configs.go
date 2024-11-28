package configs

import (
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/marine-br/golib-utils/utils/env_validator"
)

type Config struct {
	Port        string
	Environment string
	LogLevel    string
	JWT         struct {
		SecretKey string
	}
	Consul struct {
		Address     string
		ServiceName string
		ServiceHost string
		ServicePort int
	}
	RateLimit struct {
		Enabled bool
		Limit   int
		Window  time.Duration
	}
}

func LoadConfig() *Config {
	validator := env_validator.NewEnvValidator()

	cfg := &Config{
		Port:        validator.Default("APP_PORT", "8080"),
		Environment: validator.Default("APP_ENV", "development"),
		LogLevel:    validator.Default("LOG_LEVEL", "info"),
		JWT: struct {
			SecretKey string
		}{
			SecretKey: validator.Exists("JWT_SECRET_KEY"),
		},
		Consul: struct {
			Address     string
			ServiceName string
			ServiceHost string
			ServicePort int
		}{
			Address: validator.Exists("CONSUL_ADDRESS"),
		},
		RateLimit: struct {
			Enabled bool
			Limit   int
			Window  time.Duration
		}{
			Enabled: validator.DefaultBool("RATE_LIMIT_ENABLED", true),
			Limit:   validator.DefaultInt("RATE_LIMIT_PER_WINDOW", 100),
			Window:  time.Duration(validator.DefaultInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
		},
	}

	if !validator.Validate() {
		panic(validator.Errors())
	}

	return cfg
}
