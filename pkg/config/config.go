package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	APIHost string
	APIPort int

	PostgresHost        string
	PostgresPort        int
	PostgresUser        string
	PostgresPassword    string
	PostgresDB          string
	PostgresSSLMode     string
	PostgresRootCertLoc string

	JWTSecret string
}

func (c *Config) APIAddr() string { return fmt.Sprintf("%s:%d", c.APIHost, c.APIPort) }
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		// Added single quotes to accomodate empty values.
		// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
		"host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode='%s'",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSSLMode,
	)
}

func Load() (*Config, error) {
	// Load optional ".env" file for development.
	godotenv.Load()

	env := &envHelper{}
	config := &Config{
		APIHost:          env.Optional("API_HOST", ""),
		APIPort:          env.OptionalInt("API_PORT", 8080),
		PostgresHost:     env.Mandatory("POSTGRES_HOST"),
		PostgresPort:     env.MandatoryInt("POSTGRES_PORT"),
		PostgresUser:     env.Mandatory("POSTGRES_USER"),
		PostgresPassword: env.Mandatory("POSTGRES_PASSWORD"),
		PostgresDB:       env.Mandatory("POSTGRES_DB"),
		PostgresSSLMode:  env.Optional("POSTGRES_SSL_MODE", "disable"),
		JWTSecret:        env.Mandatory("JWT_SECRET"),
	}

	// Check for missing or malformed envs.
	if err := env.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
