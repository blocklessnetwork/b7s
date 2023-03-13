package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

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

func run() int {

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Initialize logging.
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	// Parse CLI flags and validate that the configuration is valid.
	cfg := parseFlags()

	// Set log level.
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Error().Err(err).Str("level", cfg.Log.Level).Msg("could not parse log level")
		return failure
	}
	log = log.Level(level)

	// Determine node role.
	role, err := parseNodeRole(cfg.Role)
	if err != nil {
		log.Error().Err(err).Str("role", cfg.Role).Msg("invalid node role specified")
		return failure
	}

	// Open the pebble database.
	pdb, err := pebble.Open(cfg.DatabasePath, &pebble.Options{})
	if err != nil {
		log.Error().Err(err).Str("db", cfg.DatabasePath).Msg("could not open pebble database")
		return failure
	}
	defer pdb.Close()

	// Create a new store.
	store := store.New(pdb)

	peerstore := peerstore.New(store)

	// Get the list of dial back peers.
	peers, err := peerstore.Peers()
	if err != nil {
		log.Error().Err(err).Msg("could not get list of dial-back peers")
		return failure
	}
	peerAddrs, err := getPeerAddresses(peers)
	if err != nil {
		log.Error().Err(err).Msg("could not get peer addresses")
		return failure
	}

	// Get the list of boot nodes addresses.
	bootNodeAddrs, err := getBootNodeAddresses(cfg.BootNodes)
	if err != nil {
		log.Error().Err(err).Msg("could not get boot node addresses")
		return failure
	}

	// Create libp2p host.
	host, err := host.New(log, cfg.Host.Address, cfg.Host.Port,
		host.WithPrivateKey(cfg.Host.PrivateKey),
		host.WithBootNodes(bootNodeAddrs),
		host.WithDialBackPeers(peerAddrs),
	)
	if err != nil {
		log.Error().Err(err).Str("key", cfg.Host.PrivateKey).Msg("could not create host")
		return failure
	}

	log.Info().
		Str("id", host.ID().String()).
		Strs("addresses", host.Addresses()).
		Int("boot_nodes", len(bootNodeAddrs)).
		Int("dial_back_peers", len(peerAddrs)).
		Msg("created host")

	// Set node options.
	opts := []node.Option{
		node.WithRole(role),
		node.WithConcurrency(cfg.Concurrency),
	}

	// If this is a worker node, initialize an executor.
	if role == blockless.WorkerNode {

		// Crete an executor.
		executor, err := executor.New(log,
			executor.WithWorkDir(cfg.Workspace),
			executor.WithRuntimeDir(cfg.Runtime),
		)
		if err != nil {
			log.Error().
				Err(err).
				Str("workspace", cfg.Workspace).
				Str("runtime", cfg.Runtime).
				Msg("could not create an executor")
			return failure
		}

		opts = append(opts, node.WithExecutor(executor))
	}

	// Create function store.
	functionStore := function.NewHandler(log, store, cfg.Workspace)

	// Instantiate node.
	node, err := node.New(log, host, store, peerstore, functionStore, opts...)
	if err != nil {
		log.Error().Err(err).Msg("could not create node")
		return failure
	}

	// Create the main context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	failed := make(chan struct{})

	// Start node main loop in a separate goroutine.
	go func() {

		log.Info().
			Str("role", role.String()).
			Msg("Blockless Node starting")

		err := node.Run(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Blockless Node failed")
			close(failed)
		} else {
			close(done)
		}

		log.Info().Msg("Blockless Node stopped")
	}()

	// If we're a head node - start the REST API.
	if role == blockless.HeadNode {

		if cfg.API == "" {
			log.Error().Err(err).Msg("REST API address is required")
			return failure
		}

		// Create echo server and iniialize logging.
		server := echo.New()
		server.HideBanner = true
		server.HidePort = true

		elog := lecho.From(log)
		server.Logger = elog
		server.Use(lecho.Middleware(lecho.Config{Logger: elog}))

		// Create an API handler.
		api := api.New(log, node)

		// Set endpoint handlers.
		server.POST("/api/v1/functions/execute", api.Execute)
		server.GET("/api/v1/functions/:id/install", api.Install)
		server.GET("/api/v1/functions/requests/:id/result", api.ExecutionResult)

		// Start API in a separate goroutine.
		go func() {

			log.Info().Msg("Node API starting")
			err := server.Start(cfg.API)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Warn().Err(err).Msg("Node API failed")
				close(failed)
			} else {
				close(done)
			}

			log.Info().Msg("Node API stopped")
		}()
	}

	select {
	case <-sig:
		log.Info().Msg("Blockless Node stopping")
	case <-done:
		log.Info().Msg("Blockless Node done")
	case <-failed:
		log.Info().Msg("Blockless Node aborted")
		return failure
	}

	// If we receive a second interrupt signal, exit immediately.
	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	return success
}
