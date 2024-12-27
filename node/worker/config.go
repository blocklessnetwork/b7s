package worker

import (
	"errors"
	"path/filepath"

	"github.com/hashicorp/go-multierror"

	"github.com/blocklessnetwork/b7s/metadata"
)

// Option can be used to set Node configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the node.
var DefaultConfig = Config{
	LoadAttributes:   DefaultAttributeLoadingSetting,
	MetadataProvider: metadata.NewNoopProvider(),
}

// Config represents the Node configuration.
type Config struct {
	Workspace        string            // Directory where we can store files needed for execution.
	LoadAttributes   bool              // Node should try to load its attributes from IPFS.
	MetadataProvider metadata.Provider // Metadata provider for the node
}

// Validate checks if the given configuration is correct.
func (c Config) Valid() error {

	var err *multierror.Error

	if !filepath.IsAbs(c.Workspace) {
		err = multierror.Append(err, errors.New("workspace must be an absolute path"))
	}

	return err.ErrorOrNil()
}

// Workspace specifies the workspace that the node can use for file storage.
func Workspace(path string) Option {
	return func(cfg *Config) {
		cfg.Workspace = path
	}
}

// AttributeLoading specifies whether node should try to load its attributes data from IPFS.
func AttributeLoading(b bool) Option {
	return func(cfg *Config) {
		cfg.LoadAttributes = b
	}
}

// MetadataProvider sets the metadata provider for the node.
func MetadataProvider(p metadata.Provider) Option {
	return func(cfg *Config) {
		cfg.MetadataProvider = p
	}
}
