package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// LoadYAML unmarshals YAML file into struct with `yaml` tags.
// It supports environment variables.
func LoadYAML(file string, dst any) error {
	f, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("file: %v", err)
	}

	// Load optional .env file.
	godotenv.Load()
	// Expand all environment variables on the file.
	f = []byte(os.ExpandEnv(string(f)))

	if err = yaml.Unmarshal(f, dst); err != nil {
		return fmt.Errorf("decode: %v", err)
	}
	return nil
}
