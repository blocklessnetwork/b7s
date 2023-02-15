package main

import (
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"

	"github.com/blocklessnetworking/b7s/api"
	"github.com/blocklessnetworking/b7s/executor"
	"github.com/blocklessnetworking/b7s/function"
	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
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

	// Initialize logging.
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	// Parse CLI flags and validate that the configuration is valid.
	cfg := parseFlags()
	err := cfg.Valid()
	if err != nil {
		log.Error().Err(err).Msg("invalid configuration")
		return failure
	}

	// Set log level.
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Error().Err(err).Str("level", cfg.Log.Level).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	// Create libp2p host.
	host, err := host.New(log, cfg.Host.Address, cfg.Host.Port, host.WithPrivateKey(cfg.Host.PrivateKey))
	if err != nil {
		log.Error().Err(err).Str("key", cfg.Host.PrivateKey).Msg("could not create host")
		return failure
	}

	hostIDs := host.IDs()
	log.Info().Strs("ids", hostIDs).Msg("created host")

	// Open the pebble database.
	pdb, err := pebble.Open(cfg.DatabasePath, &pebble.Options{})
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

	// Set node options.
	opts := []node.Option{
		node.WithRole(role),
	}

	// If this is a worker node, initialize an executor.
	if role == blockless.WorkerNode {

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

		opts = append(opts, node.WithExecute(executor))
	}

	// Create function handler.
	functionHandler := function.New(log, store, cfg.Workspace)

	// Instantiate node.
	node, err := node.New(log, host, store, peerstore, functionHandler, opts...)
	if err != nil {
		log.Error().Err(err).Msg("could not create node")
		return failure
	}

	// If we're a head node - start the REST API.
	if role == blockless.HeadNode {

		// Create echo server and iniialize logging.
		server := echo.New()
		server.HideBanner = true
		server.HidePort = true

		elog := lecho.From(log)
		server.Logger = elog
		server.Use(lecho.Middleware(lecho.Config{Logger: elog}))

		// Create an API handler.
		api := api.New(node)

		// Set endpoint handlers.
		server.POST("/api/v1/functions/execute", api.Execute)
		server.GET("/api/v1/functions/:id/install", api.Install)
		server.GET("/api/v1/functions/requests/:id/result", api.ExecutionResult)
	}

	return failure
}
