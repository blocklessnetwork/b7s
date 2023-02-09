package host

// defaultConfig used to create Host.
var defaultConfig = Config{
	PrivateKey: "",
}

// Config represents the Host configuration.
type Config struct {
	PrivateKey string
}

// WithPrivateKey specifies the private key for the Host.
func WithPrivateKey(filepath string) func(*Config) {
	return func(cfg *Config) {
		cfg.PrivateKey = filepath
	}
}
