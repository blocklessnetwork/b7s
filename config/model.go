package config

// Config describes the Blockless configuration options.
type Config struct {
	Log                  Log
	PeerDatabasePath     string
	FunctionDatabasePath string
	Role                 string
	BootNodes            []string
	Concurrency          uint

	Host    Host
	API     string
	Runtime string

	CPUPercentage float64
	MemoryMaxKB   int64

	Workspace string
}

// Host describes the libp2p host that the node will use.
type Host struct {
	Port            uint
	Address         string
	PrivateKey      string
	DialBackPort    uint
	DialBackAddress string
	Websocket       bool
}

// Log describes the logging configuration.
type Log struct {
	Level string
}
