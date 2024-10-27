package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type envHelper struct {
	missing   []string // List of mandatory environment variables that are missing.
	malformed []string // List of environment variables that failed parsing.
}

func (env *envHelper) Validate() error {
	if len(env.missing) > 0 {
		return fmt.Errorf("missing mandatory env: %s", strings.Join(env.missing, ", "))
	}
	if len(env.malformed) > 0 {
		return fmt.Errorf("malformed env values: %s", strings.Join(env.malformed, ", "))
	}
	return nil
}

func (env *envHelper) Mandatory(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
	}
	return value
}

func (env *envHelper) Optional(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func (env *envHelper) MandatoryInt(key string) int {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
		return 0
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return 0
	}
	return value
}

func (env *envHelper) OptionalInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return fallback
	}
	return value
}

func (env *envHelper) MandatoryBool(key string) bool {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
		return false
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return false
	}
	return value
}

func (env *envHelper) OptionalBool(key string, fallback bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return fallback
	}
	return value
}

func (env *envHelper) MandatoryDuration(key string) time.Duration {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		env.missing = append(env.missing, key)
		return 0
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return 0
	}
	return value
}

func (env *envHelper) OptionalDuration(key string, fallback time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		env.malformed = append(env.malformed, key)
		return fallback
	}
	return value
}
