package configs

import (
	"fmt"
	"strconv"

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
}

func LoadConfig() *Config {
	validator := env_validator.NewEnvValidator()
	servicePort, err := strconv.Atoi(validator.Default("CONSUL_SERVICE_PORT", "8080"))

	if err != nil {
		panic(fmt.Errorf("failed to convert CONSUL_SERVICE_PORT to int: %w", err))
	}

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
			Address:     validator.Exists("CONSUL_ADDRESS"),
			ServiceName: validator.Exists("CONSUL_SERVICE_NAME"),
			ServiceHost: validator.Exists("CONSUL_SERVICE_HOST"),
			ServicePort: servicePort,
		},
	}

	if !validator.Validate() {
		panic(validator.Errors())
	}

	return cfg
}
