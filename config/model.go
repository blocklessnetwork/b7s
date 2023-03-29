package config

import (
	"time"
)

// Config describes the Blockless configuration options.
type Config struct {
	Log          Log
	DatabasePath string
	Role         string
	BootNodes    []string
	Concurrency  uint

	Host    Host
	API     string
	Runtime string

	CPUTime     time.Duration
	MemoryMaxKB int64

	Workspace string
}

// Host describes the libp2p host that the node will use.
type Host struct {
	Port            uint
	Address         string
	PrivateKey      string
	DialBackPort    uint
	DialBackAddress string
}

// Log describes the logging configuration.
type Log struct {
	Level string
}
