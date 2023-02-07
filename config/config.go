package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

// Load will load the config file from the given location.
func Load(file string) (*Config, error) {

	// Read config file.
	payload, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	// Unmarshal file.
	var config Config
	err = yaml.Unmarshal(payload, &config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal file: %w", err)
	}

	// Validate configuration.
	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &config, nil
}
