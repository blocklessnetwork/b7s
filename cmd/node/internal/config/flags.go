package config

import (
	"github.com/blocklessnetwork/b7s/node"
)

// Default values.
const (
	DefaultPort         = 0
	DefaultAddress      = "0.0.0.0"
	DefaultRole         = "worker"
	DefaultPeerDB       = "peer-db"
	DefaultFunctionDB   = "function-db"
	DefaultConcurrency  = uint(node.DefaultConcurrency)
	DefaultUseWebsocket = false
	DefaultWorkspace    = ""
)

// CLI flag names
const (
	// General
	flagConfig      = "config"
	flagRole        = "role"
	flagConcurrency = "concurrency"
	flagBootNodes   = "boot-nodes"
	flagWorkspace   = "workspace"
	flagAttributes  = "attributes"
	flagPeerDB      = "peer-db"
	flagFunctionDB  = "function-db"
	flagTopics      = "topics"
	// Connectivity
	flagAddress               = "address"
	flagPort                  = "port"
	flagPrivateKey            = "private-key"
	flagDialbackAddress       = "dialback-address"
	flagDialbackPort          = "dialback-port"
	flagWebsocket             = "websocket"
	flagWebsocketPort         = "websocket-port"
	flagWebsocketDialbackPort = "websocket-dialback-port"
	// Head
	flagRestAPI = "rest-api"
	// Worker
	flagRuntimePath = "runtime-path"
	flagRuntimeCLI  = "runtime-cli"
	flagCPULimit    = "cpu-percentage-limit"
	flagMemoryLimit = "memory-limit"
	// Log
	flagLogLevel = "log-level"
)
