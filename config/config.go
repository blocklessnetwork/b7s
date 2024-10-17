package config

import (
	"github.com/blocklessnetwork/b7s/execution/limits"
	"github.com/blocklessnetwork/b7s/node"
)

// Default values.
const (
	DefaultPort         = uint(0)
	DefaultAddress      = "0.0.0.0"
	DefaultRole         = "worker"
	DefaultConcurrency  = uint(node.DefaultConcurrency)
	DefaultUseWebsocket = false
	DefaultLogLevel     = "info"
)

// Default names for storage directories.
const (
	DefaultDBName        = "db"
	DefaultWorkspaceName = "workspace"
)

var DefaultConfig = Config{
	Role:        DefaultRole,
	Concurrency: DefaultConcurrency,
	Log: Log{
		Level: DefaultLogLevel,
	},
	Connectivity: Connectivity{
		Address:   DefaultAddress,
		Port:      DefaultPort,
		Websocket: DefaultUseWebsocket,
	},
	Worker: Worker{
		CgroupMountpoint:   limits.DefaultMountpoint,
		CgroupName:         "b7s",
		CPUPercentageLimit: 1.0,
		MemoryLimitKB:      -1,
	},
}

// Config describes the Blockless configuration options.
// NOTE: DO NOT use TABS in struct tags - spaces only!
// NOTE: When adding CLI flags (using the `flag` struct tag) - add the description for (for the flag long version, not the shorthand) it in getFlagDescription() below.
type Config struct {
	Role           string   `koanf:"role"            flag:"role,r"`
	Concurrency    uint     `koanf:"concurrency"     flag:"concurrency,c"`
	BootNodes      []string `koanf:"boot-nodes"      flag:"boot-nodes"`
	Workspace      string   `koanf:"workspace"       flag:"workspace"`       // TODO: Check - does a head node ever use a workspace?
	LoadAttributes bool     `koanf:"load-attributes" flag:"load-attributes"` // TODO: Head node probably doesn't need attributes..?
	Topics         []string `koanf:"topics"          flag:"topics"`

	DB string `koanf:"db" flag:"db"`

	Log          Log          `koanf:"log"`
	Connectivity Connectivity `koanf:"connectivity"`
	Head         Head         `koanf:"head"`
	Worker       Worker       `koanf:"worker"`
	Telemetry    Telemetry    `koanf:"telemetry"`
}

// Log describes the logging configuration.
type Log struct {
	Level string `koanf:"level" flag:"log-level,l"`
}

type Head struct {
	RestAPI string `koanf:"rest-api" flag:"rest-api"`
}

// ConfigOptionInfo describes a specific configuration option, it's location in the config file and
// corresponding CLI flags and environment variables. It can be used to generate documentation for the b7s node.
type ConfigOptionInfo struct {
	Name     string         `json:"name,omitempty"      yaml:"name,omitempty"`
	FullPath string         `json:"full_path,omitempty" yaml:"full_path,omitempty"`
	CLI      CLIFlag        `json:"cli,omitempty"       yaml:"cli,omitempty"`
	Env      string         `json:"env-var,omitempty"   yaml:"env-var,omitempty"`
	Children []ConfigOption `json:"children,omitempty"  yaml:"children,omitempty"`
	Type     string         `json:"type,omitempty"      yaml:"type,omitempty"`
}

func getFlagDescription(flag string) string {

	switch flag {
	case "role":
		return "role this node will have in the Blockless protocol (head or worker)"
	case "concurrency":
		return "maximum number of requests node will process in parallel"
	case "boot-nodes":
		return "list of addresses that this node will connect to on startup, in multiaddr format"
	case "workspace":
		return "directory that the node can use for file storage"
	case "load-attributes":
		return "node should try to load its attribute data from IPFS"
	case "topics":
		return "topics node should subscribe to"
	case "db":
		return "path to the database used for persisting peer and function data"
	case "log-level":
		return "log level to use"
	case "address":
		return "address that the b7s host will use"
	case "port":
		return "port that the b7s host will use"
	case "private-key":
		return "private key that the b7s host will use"
	case "websocket":
		return "should the node use websocket protocol for communication"
	case "dialback-address":
		return "external address that the b7s host will advertise"
	case "dialback-port":
		return "external port that the b7s host will advertise"
	case "websocket-port":
		return "port to use for websocket connections"
	case "websocket-dialback-port":
		return "external port that the b7s host will advertise for websocket connections"
	case "connection-count":
		return "maximum number of connections the b7s host will aim to have"
	case "rest-api":
		return "address where the head node REST API will listen on"
	case "runtime-path":
		return "Blockless Runtime location (used by the worker node)"
	case "runtime-cli":
		return "runtime CLI name (used by the worker node)"
	case "cpu-percentage-limit":
		return "amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited"
	case "memory-limit":
		return "memory limit (kB) for Blockless Functions"
	case "no-dialback-peers":
		return "start without dialing back peers from previous runs"
	case "must-reach-boot-nodes":
		return "halt node if we fail to reach boot nodes on start"
	case "disable-connection-limits":
		return "disable libp2p connection limits (experimental)"
	case "enable-tracing":
		return "emit tracing data"
	case "enable-metrics":
		return "emit metrics"
	case "tracing-grpc-endpoint":
		return "tracing exporter GRPC endpoint"
	case "tracing-http-endpoint":
		return "tracing exporter HTTP endpoint"
	case "prometheus-address":
		return "address where prometheus metrics will be served"
	default:
		return ""
	}
}
