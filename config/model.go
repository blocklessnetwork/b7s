package config

// Config describes the Blockless configuration options.
type Config struct {
	Log          Log
	DatabasePath string
	Role         string
	BootNodes    []string

	Host    Host
	API     string
	Runtime string

	Workspace string
}

// Host describes the libp2p host that the node will use.
type Host struct {
	Port       uint
	Address    string
	PrivateKey string
}

// Log describes the logging configuration.
type Log struct {
	Level string
}
