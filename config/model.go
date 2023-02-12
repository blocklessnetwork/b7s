package config

import (
	"github.com/go-playground/validator/v10"
)

// Config describes the Blockless configuration options.
type Config struct {
	Log          Log    `validate:"required"`
	DatabasePath string `validate:"required"`
	Role         string `validate:"oneof=head worker"`
	BootNodes    []string

	Host    Host   `validate:"required"`
	API     string `validate:"required_if=role head,excluded_if=role worker"`
	Runtime string `validate:"dir,required_if=role worker,excluded_if=role head"`

	Workspace string `validate:"required"`
}

// Host describes the libp2p host that the node will use.
type Host struct {
	Port       uint
	Address    string `validate:"required"`
	PrivateKey string `validate:"omitempty,file"`
}

// Log describes the logging configuration.
type Log struct {
	Level string `validate:"required"`
}

// Valid will check if the provided configuration is valid, and return an error if not.
func (c Config) Valid() error {
	return validator.New().Struct(c)
}
