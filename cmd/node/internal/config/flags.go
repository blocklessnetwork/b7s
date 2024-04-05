package config

import (
	"github.com/blocklessnetwork/b7s/node"
	"github.com/spf13/pflag"
)

// Default values.
const (
	DefaultPort         = uint(0)
	DefaultAddress      = "0.0.0.0"
	DefaultRole         = "worker"
	DefaultPeerDB       = "peer-db"
	DefaultFunctionDB   = "function-db"
	DefaultConcurrency  = uint(node.DefaultConcurrency)
	DefaultUseWebsocket = false
	DefaultWorkspace    = "workspace"
)

type configOption struct {
	flag  string      // long flag name - should be the same as the `koanf` tag in the Config type.
	short string      // shorthand - single letter alternative to the long flag name
	group configGroup // group - defined in which section of the config file this option lives.
	usage string      // description
}

// Config options.
var (
	// Root group.
	roleCfg = configOption{
		flag:  "role",
		short: "r",
		group: rootGroup,
		usage: "role this note will have in the Blockless protocol (head or worker)",
	}
	concurrencyCfg = configOption{
		flag:  "concurrency",
		short: "c",
		group: rootGroup,
		usage: "maximum number of requests node will process in parallel",
	}
	bootNodesCfg = configOption{
		flag:  "boot-nodes",
		group: rootGroup,
		usage: "list of addresses that this node will connect to on startup, in multiaddr format",
	}
	workspaceCfg = configOption{
		flag:  "workspace",
		group: rootGroup,
		usage: "directory that the node can use for file storage",
	}
	attributesCfg = configOption{
		flag:  "attributes",
		group: rootGroup,
		usage: "node should try to load its attribute data from IPFS",
	}
	peerDBCfg = configOption{
		flag:  "peer-db",
		group: rootGroup,
		usage: "path to the database used for persisting peer data",
	}
	functionDBCfg = configOption{
		flag:  "function-db",
		group: rootGroup,
		usage: "path to the database used for persisting function data",
	}
	topicsCfg = configOption{
		flag:  "topics",
		group: rootGroup,
		usage: "topics node should subscribe to",
	}

	// Log group.
	logLevelCfg = configOption{
		flag:  "log-level",
		short: "l",
		group: logGroup,
		usage: "log level to use",
	}

	// Connectivity group.
	addressCfg = configOption{
		flag:  "address",
		short: "a",
		group: connectivityGroup,
		usage: "address that the b7s host will use",
	}
	portCfg = configOption{
		flag:  "port",
		short: "p",
		group: connectivityGroup,
		usage: "port that the b7s host will use",
	}
	privateKeyCfg = configOption{
		flag:  "private-key",
		group: connectivityGroup,
		usage: "private key that the b7s host will use",
	}
	websocketCfg = configOption{
		flag:  "websocket",
		short: "w",
		group: connectivityGroup,
		usage: "should the node use websocket protocol for communication",
	}
	websocketPortCfg = configOption{
		flag:  "websocket-port",
		group: connectivityGroup,
		usage: "port to use for websocket connections",
	}
	dialbackAddressCfg = configOption{
		flag:  "dialback-address",
		group: connectivityGroup,
		usage: "external address that the b7s host will advertise",
	}
	dialbackPortCfg = configOption{
		flag:  "dialback-port",
		group: connectivityGroup,
		usage: "external port that the b7s host will advertise",
	}
	websocketDialbackPortCfg = configOption{
		flag:  "websocket-dialback-port",
		group: connectivityGroup,
		usage: "external port that the b7s host will advertise for websocket connections",
	}

	// Worker flags.
	runtimePathCfg = configOption{
		flag:  "runtime-path",
		group: workerGroup,
		usage: "Blockless Runtime location (used by the worker node)",
	}
	runtimeCLICfg = configOption{
		flag:  "runtime-cli",
		group: workerGroup,
		usage: "runtime CLI name (used by the worker node)",
	}
	cpuLimitCfg = configOption{
		flag:  "cpu-percentage-limit",
		group: workerGroup,
		usage: "amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited",
	}
	memLimitCfg = configOption{
		flag:  "memory-limit",
		group: workerGroup,
		usage: "memory limit (kB) for Blockless Functions",
	}

	// Head node flags.
	restAPICfg = configOption{
		flag:  "rest-api",
		group: headGroup,
		usage: "address where the head node REST API will listen on",
	}
)

// This helper type is a thin wrapper around the pflag.FlagSet.
// Added functionality is the accounting of added flags.
// This is needed/useful when we're translating flags between the structured format (yaml file) and the flat structure (CLI flags).
type cliFlags struct {
	fs      *pflag.FlagSet
	options []configOption
}

func newCliFlags() *cliFlags {

	fs := pflag.NewFlagSet("b7s-node", pflag.ExitOnError)
	fs.SortFlags = false

	return &cliFlags{
		fs:      fs,
		options: make([]configOption, 0),
	}
}

func (c *cliFlags) stringFlag(cfg configOption, defaultValue string) {
	c.fs.StringP(cfg.flag, cfg.short, defaultValue, cfg.usage)
	c.options = append(c.options, cfg)
}

func (c *cliFlags) boolFlag(cfg configOption, defaultValue bool) {
	c.fs.BoolP(cfg.flag, cfg.short, defaultValue, cfg.usage)
	c.options = append(c.options, cfg)
}

func (c *cliFlags) uintFlag(cfg configOption, defaultValue uint) {
	c.fs.UintP(cfg.flag, cfg.short, defaultValue, cfg.usage)
	c.options = append(c.options, cfg)
}

func (c *cliFlags) int64Flag(cfg configOption, defaultValue int64) {
	c.fs.Int64P(cfg.flag, cfg.short, defaultValue, cfg.usage)
	c.options = append(c.options, cfg)
}

func (c *cliFlags) float64Flag(cfg configOption, defaultValue float64) {
	c.fs.Float64P(cfg.flag, cfg.short, defaultValue, cfg.usage)
	c.options = append(c.options, cfg)
}

func (c *cliFlags) stringSliceFlag(cfg configOption, defaultValue []string) {
	c.fs.StringSliceP(cfg.flag, cfg.short, defaultValue, cfg.usage)
	c.options = append(c.options, cfg)
}

func (c *cliFlags) groups() map[string]configGroup {

	groups := make(map[string]configGroup)
	for _, option := range c.options {
		groups[option.flag] = option.group
	}

	return groups
}
