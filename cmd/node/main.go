package main

import (
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/config"
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

// TODO: Logging format - JSON vs text.
// TODO: Two variants for config loading - look for config file in CWD or explicitely from the flag value.

func run() int {

	var (
		flagAddress  string
		flagDB       string
		flagConfig   string
		flagLogLevel string
		flagPort     uint

		flagNodeRole   string
		flagPrivateKey string
	)

	pflag.StringVarP(&flagAddress, "address", "a", defaultAddress, "address to use")
	pflag.StringVarP(&flagDB, "db-path", "d", defaultDB, "path to the node database")
	pflag.StringVarP(&flagConfig, "config", "c", "config.yaml", "path to config file")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level to use")
	pflag.UintVarP(&flagPort, "port", "p", defaultPort, "port number to use - random port if 0")

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
	host, err := host.New(flagAddress, flagPort, host.WithPrivateKey(flagPrivateKey))
	if err != nil {
		log.Error().Err(err).Str("key", flagPrivateKey).Msg("could not create host")
		return failure
	}

	hostIDs := host.IDs()
	log.Info().Strs("ids", hostIDs).Msg("created host")

	// TODO: Implement messaging.ListenMessages functionality from old host package.

	// TODO: If we're listening on 0.0.0.0 we'll have multiple IDs - one for each network interface.
	// It may still make sense to use the /ip4/0.0.0.0/tcp/<port>/p2p/<host-id>_appDB for the DB - instead of using multiple ones.
	// But also - do we even need some kind of easily switchable databases? I assume we'll typically keep one. If someone wants to switch,
	// they can point the executable to a different DB.

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

	// Instantiate node.
	node, err := node.New(log, host, peerstore, node.WithRole(role))
	if err != nil {
		log.Error().Err(err).Msg("could not create node")
		return failure
	}

	// TODO: Remove
	_ = node

	return failure
}
