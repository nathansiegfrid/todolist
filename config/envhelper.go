package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type EnvHelper struct {
	missing   []string // List of mandatory environment variables that are missing.
	malformed []string // List of environment variables that failed parsing.
}

func (env *EnvHelper) Validate() error {
	if len(env.missing) > 0 {
		return fmt.Errorf("missing mandatory env: %s", strings.Join(env.missing, ", "))
	}
	if len(env.malformed) > 0 {
		return fmt.Errorf("malformed env values: %s", strings.Join(env.malformed, ", "))
	}
	return nil
}

func (env *EnvHelper) Mandatory(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
	}

	return value
}

func (env *EnvHelper) MandatoryInt(key string) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
		return 0
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return 0
	}

	return valueInt
}

func (env *EnvHelper) MandatoryDuration(key string) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
		return 0
	}

	valueDuration, err := time.ParseDuration(value)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return 0
	}

	return valueDuration
}

func (env *EnvHelper) MandatoryBool(key string) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
		return false
	}

	valueBool, err := strconv.ParseBool(value)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return false
	}

	return valueBool
}

func (env *EnvHelper) Optional(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func (env *EnvHelper) OptionalInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return fallback
	}

	return valueInt
}

func (env *EnvHelper) OptionalBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueBool, err := strconv.ParseBool(value)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return fallback
	}

	return valueBool
}

func (env *EnvHelper) OptionalDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueDuration, err := time.ParseDuration(value)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return fallback
	}

	return valueDuration
}
