package main

import (
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog"

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

func main() {
	os.Exit(run())
}

// TODO: RestAPI za node i execution memstore

func run() int {

	// Set log level.
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	// Parse CLI flags.
	cfg := parseFlags()

	err := cfg.Valid()
	if err != nil {
		log.Error().Err(err).Msg("invalid configuration")
		return failure
	}

	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Error().Err(err).Str("level", cfg.Log.Level).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	// Create host.
	host, err := host.New(log, cfg.Host.Address, cfg.Host.Port, host.WithPrivateKey(cfg.Host.PrivateKey))
	if err != nil {
		log.Error().Err(err).Str("key", cfg.Host.PrivateKey).Msg("could not create host")
		return failure
	}

	hostIDs := host.IDs()
	log.Info().Strs("ids", hostIDs).Msg("created host")

	// Open the pebble database.
	opts := pebble.Options{}
	pdb, err := pebble.Open(cfg.DatabasePath, &opts)
	if err != nil {
		log.Error().Err(err).Str("db", cfg.DatabasePath).Msg("could not open pebble database")
		return failure
	}
	defer pdb.Close()

	// Create a new store.
	store := store.New(pdb)

	// Determine node role.
	role, err := parseNodeRole(cfg.Role)
	if err != nil {
		log.Error().Err(err).Str("role", cfg.Role).Msg("invalid node role specified")
		return failure
	}

	log.Info().Str("role", role.String()).Msg("starting node")

	peerstore := peerstore.New(store)

	// Crete an executor.
	executor, err := executor.New(log, cfg.Workspace, cfg.Runtime)
	if err != nil {
		log.Error().
			Err(err).
			Str("workspace", cfg.Workspace).
			Str("runtime", cfg.Runtime).
			Msg("could not create an executor")
		return failure
	}

	// Create function handler.
	functionHandler := function.New(log, store, cfg.Workspace)

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
