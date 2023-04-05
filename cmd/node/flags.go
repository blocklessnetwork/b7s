package main

import (
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/config"
	"github.com/blocklessnetworking/b7s/node"
)

// Default values.
const (
	defaultPort        = 0
	defaultAddress     = "0.0.0.0"
	defaultPeerDB      = "peer-db"
	defaultFunctionDB  = "function-db"
	defaultConcurrency = uint(node.DefaultConcurrency)

	defaultRole = "worker"
)

func parseFlags() *config.Config {

	var cfg config.Config

	pflag.StringVarP(&cfg.Log.Level, "log-level", "l", "info", "log level to use")
	pflag.StringVarP(&cfg.PeerDatabasePath, "peer-db", "d", defaultPeerDB, "path to the database used for persisting peer data")
	pflag.StringVarP(&cfg.FunctionDatabasePath, "function-db", "d", defaultFunctionDB, "path to the database used for persisting function data")

	// Node configuration.
	pflag.StringVarP(&cfg.Role, "role", "r", defaultRole, "role this note will have in the Blockless protocol (head or worker)")
	pflag.StringVarP(&cfg.Host.Address, "address", "a", defaultAddress, "address that the b7s host will use")
	pflag.UintVarP(&cfg.Host.Port, "port", "p", defaultPort, "port that the b7s host will use")

	pflag.StringVarP(&cfg.Host.DialBackAddress, "dialback-address", "", defaultAddress, "external address that the b7s host will advertise")
	pflag.UintVarP(&cfg.Host.DialBackPort, "dialback-port", "", defaultPort, "external port that the b7s host will advertise")

	pflag.StringVar(&cfg.Host.PrivateKey, "private-key", "", "private key that the b7s host will use")
	pflag.UintVarP(&cfg.Concurrency, "concurrency", "c", defaultConcurrency, "maximum number of requests node will process in parallel")
	pflag.StringVar(&cfg.API, "rest-api", "", "address where the head node REST API will listen on")
	pflag.StringSliceVar(&cfg.BootNodes, "boot-nodes", nil, "list of addresses that this node will connect to on startup, in multiaddr format")

	pflag.StringVar(&cfg.Workspace, "workspace", "./workspace", "directory that the node can use for file storage")
	pflag.StringVar(&cfg.Runtime, "runtime", "", "runtime address (used by the worker node)")

	pflag.Float64Var(&cfg.CPUPercentage, "cpu-percentage-limit", 1.0, "amount of CPU time allowed for Blockless Functions in the 0-1 range, 1 being unlimited")
	pflag.Int64Var(&cfg.MemoryMaxKB, "memory-limit", 0, "memory limit (kB) for Blockless Functions")

	pflag.CommandLine.SortFlags = false

	pflag.Parse()

	return &cfg
}
