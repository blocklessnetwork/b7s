package main

import (
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/config"
	"github.com/blocklessnetworking/b7s/executor"
	"github.com/blocklessnetworking/b7s/function"
	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/node"
	"github.com/blocklessnetworking/b7s/peerstore"
	"github.com/blocklessnetworking/b7s/store"
)

const (
	success = 0
	failure = 1
)

const (
	// TODO: Default port for head node is 9527? Move to config if so.
	defaultPort    = 0
	defaultAddress = "0.0.0.0"
	defaultDB      = "db"
)

func main() {
	os.Exit(run())
}

// TODO: Have flags set as part of a struct; then load them from config and override via CLI!
// TODO: Workspace and runtime directories are overkill for the CLI flags.
// TODO: RestAPI za node i execution memstore

func run() int {

	var (
		flagAddress   string
		flagDB        string
		flagConfig    string
		flagLogLevel  string
		flagPort      uint
		flagRuntime   string
		flagWorkspace string

		flagNodeRole   string
		flagPrivateKey string
	)

	pflag.StringVarP(&flagAddress, "address", "a", defaultAddress, "address to use")
	pflag.StringVarP(&flagDB, "db-path", "d", defaultDB, "path to the node database")
	pflag.StringVarP(&flagConfig, "config", "c", "config.yaml", "path to config file")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level to use")
	pflag.UintVarP(&flagPort, "port", "p", defaultPort, "port number to use - random port if 0")
	pflag.StringVarP(&flagRuntime, "runtime-dir", "w", ".", "runtime directory where blockless-cli can be found")
	pflag.StringVarP(&flagWorkspace, "workspace-dir", "w", ".", "workspace directory")

	pflag.StringVar(&flagNodeRole, "node-role", "", "node role (head or worker)")
	pflag.StringVar(&flagPrivateKey, "private-key", "", "private key to use")

	pflag.Parse()

	// Set log level.
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Str("level", flagLogLevel).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	// Load configuration.
	cfg, err := config.Load(flagConfig)
	if err != nil {
		log.Error().Err(err).Str("config", flagConfig).Msg("could not load configuration")
		return failure
	}

	// TODO: Remove
	_ = cfg

	// Create host.
	host, err := host.New(log, flagAddress, flagPort, host.WithPrivateKey(flagPrivateKey))
	if err != nil {
		log.Error().Err(err).Str("key", flagPrivateKey).Msg("could not create host")
		return failure
	}

	hostIDs := host.IDs()
	log.Info().Strs("ids", hostIDs).Msg("created host")

	// Open the pebble database.
	opts := pebble.Options{}
	pdb, err := pebble.Open(flagDB, &opts)
	if err != nil {
		log.Error().Err(err).Str("db", flagDB).Msg("could not open pebble database")
		return failure
	}
	defer pdb.Close()

	// Create a new store.
	store, err := store.New(pdb)
	if err != nil {
		log.Error().Err(err).Str("db", flagDB).Msg("could not connect to the database")
		return failure
	}

	// Determine node role.
	role, err := parseNodeRole(flagNodeRole)
	if err != nil {
		log.Error().Err(err).Str("role", flagNodeRole).Msg("invalid node role specified")
		return failure
	}

	log.Info().Str("role", role.String()).Msg("starting node")

	peerstore := peerstore.New(store)

	// Crete an executor.
	executor, err := executor.New(log, flagWorkspace, flagRuntime)
	if err != nil {
		log.Error().
			Err(err).
			Str("workspace", flagWorkspace).
			Str("runtime", flagRuntime).
			Msg("could not create an executor")
		return failure
	}

	// Create function handler.
	functionHandler := function.New(log, store, flagWorkspace)

	// Instantiate node.
	node, err := node.New(log, host, store, executor, peerstore, functionHandler, node.WithRole(role))
	if err != nil {
		log.Error().Err(err).Msg("could not create node")
		return failure
	}

	// TODO: Remove
	_ = node

	return failure
}
