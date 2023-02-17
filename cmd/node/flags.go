package main

import (
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/config"
)

// Default values.
const (
	defaultPort    = 0
	defaultAddress = "0.0.0.0"
	defaultDB      = "db"

	defaultRole = "worker"
)

func parseFlags() *config.Config {

	var cfg config.Config

	pflag.StringVarP(&cfg.Log.Level, "log-level", "l", "info", "log level to use")
	pflag.StringVarP(&cfg.DatabasePath, "db", "d", defaultDB, "path to the database used for persisting node data")

	// Node configuration.
	pflag.StringVarP(&cfg.Role, "role", "r", defaultRole, "role this note will have in the Blockless protocol (head or worker)")
	pflag.StringVarP(&cfg.Host.Address, "address", "a", defaultAddress, "address that the libp2p host will use")
	pflag.UintVarP(&cfg.Host.Port, "port", "p", defaultPort, "port that the libp2p host will use")
	pflag.StringVar(&cfg.Host.PrivateKey, "private-key", "", "private key that the libp2p host will use")
	pflag.StringVar(&cfg.API, "rest-api", "", "address where the head node REST API will listen on")
	pflag.StringSliceVar(&cfg.BootNodes, "boot-nodes", nil, "list of addresses that this node will connect to on startup, in multiaddr format")

	pflag.StringVar(&cfg.Workspace, "workspace", "./workspace", "directory that the node can use for file storage")
	pflag.StringVar(&cfg.Runtime, "runtime", "", "runtime address (used by the worker node)")

	pflag.Parse()

	return &cfg
}
