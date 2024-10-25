package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	APIHost string
	APIPort int

	PGHost        string
	PGPort        int
	PGUser        string
	PGPassword    string
	PGDatabase    string
	PGSSLMode     string
	PGRootCertLoc string

	JWTSecret string
}

func (c *Config) APIAddr() string {
	return fmt.Sprintf("%s:%d", c.APIHost, c.APIPort)
}

func (c *Config) PGString() string {
	return fmt.Sprintf(
		// Added single quotes to accomodate empty values.
		// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
		"host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode='%s' sslrootcert='%s'",
		c.PGHost, c.PGPort, c.PGUser, c.PGPassword, c.PGDatabase, c.PGSSLMode, c.PGRootCertLoc,
	)
}

func Load() (*Config, error) {
	// Load optional ".env" file for development.
	godotenv.Load()

	env := &EnvHelper{}
	config := &Config{
		APIHost:       env.Optional("API_HOST", ""),
		APIPort:       env.OptionalInt("API_PORT", 8080),
		PGHost:        env.Mandatory("PG_HOST"),
		PGPort:        env.MandatoryInt("PG_PORT"),
		PGUser:        env.Mandatory("PG_USER"),
		PGPassword:    env.Mandatory("PG_PASSWORD"),
		PGDatabase:    env.Mandatory("PG_DATABASE"),
		PGSSLMode:     env.Optional("PG_SSL_MODE", "disable"),
		PGRootCertLoc: env.Optional("PG_ROOT_CERT_LOC", ""),
		JWTSecret:     env.Mandatory("JWT_SECRET"),
	}

	// Check for missing or malformed envs.
	if err := env.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
