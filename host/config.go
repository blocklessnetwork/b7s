package host

// TODO: Check for the 'random identity' part - is this really the case.

// defaultConfig will not use a private key path and will start
// with a random identity.
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
