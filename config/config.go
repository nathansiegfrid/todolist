package config

import "github.com/joho/godotenv"

type Config struct {
	APIHost    string
	APIPort    int
	APIVersion string

	PGHost        string
	PGPort        int
	PGUser        string
	PGPassword    string
	PGDatabase    string
	PGSSLMode     string
	PGRootCertLoc string
}

func Load() (*Config, error) {
	// Load optional ".env" file for development.
	godotenv.Load()

	env := &EnvHelper{}
	config := &Config{
		APIHost:    env.Optional("API_HOST", ""),
		APIPort:    env.OptionalInt("API_PORT", 8080),
		APIVersion: env.Optional("API_VERSION", "1.0.0"),

		PGHost:        env.Mandatory("PG_HOST"),
		PGPort:        env.MandatoryInt("PG_PORT"),
		PGUser:        env.Mandatory("PG_USER"),
		PGPassword:    env.Mandatory("PG_PASSWORD"),
		PGDatabase:    env.Mandatory("PG_DATABASE"),
		PGSSLMode:     env.Optional("PG_SSL_MODE", "disable"),
		PGRootCertLoc: env.Optional("PG_ROOT_CERT_LOC", ""),
	}

	// Check for missing or malformed envs.
	if err := env.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
