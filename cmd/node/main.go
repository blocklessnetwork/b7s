package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"

	"github.com/blocklessnetwork/b7s/api"
	"github.com/blocklessnetwork/b7s/config"
	"github.com/blocklessnetwork/b7s/executor"
	"github.com/blocklessnetwork/b7s/executor/limits"
	"github.com/blocklessnetwork/b7s/fstore"
	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node"
	"github.com/blocklessnetwork/b7s/peerstore"
	"github.com/blocklessnetwork/b7s/store"
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
	log := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	// Parse CLI flags and validate that the configuration is valid.
	cfg, err := config.Load()
	if err != nil {
		log.Error().Err(err).Msg("could not read configuration")
		return failure
	}

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

	// If we have a key, use path that corresponds to that key e.g. `.b7s_<peer-id>`.
	nodeDir := ""
	if cfg.Connectivity.PrivateKey != "" {
		id, err := peerIDFromKey(cfg.Connectivity.PrivateKey)
		if err != nil {
			log.Error().Err(err).Str("key", cfg.Connectivity.PrivateKey).Msg("could not read private key")
			return failure
		}

		nodeDir = generateNodeDirName(id)
	} else {
		nodeDir, err = os.MkdirTemp("", ".b7s_*")
		if err != nil {
			log.Error().Err(err).Msg("could not create node directory")
			return failure
		}
	}

	// Set relevant working paths for workspace, peerDB and functionDB.
	// If paths were set using the CLI flags, use those. Else, use generated path, e.g. .b7s_<peer-id>/<default-option-for-directory>.
	updateDirPaths(nodeDir, cfg)

	log.Info().
		Str("workspace", cfg.Workspace).
		Str("peer_db", cfg.PeerDB).
		Str("function_db", cfg.FunctionDB).
		Msg("filepaths used by the node")

	// Convert workspace path to an absolute one.
	workspace, err := filepath.Abs(cfg.Workspace)
	if err != nil {
		log.Error().Err(err).Str("path", cfg.Workspace).Msg("could not determine absolute path for workspace")
		return failure
	}
	cfg.Workspace = workspace

	// Open the pebble peer database.
	pdb, err := pebble.Open(cfg.PeerDB, &pebble.Options{Logger: &pebbleNoopLogger{}})
	if err != nil {
		log.Error().Err(err).Str("db", cfg.PeerDB).Msg("could not open pebble peer database")
		return failure
	}
	defer pdb.Close()

	// Create a new store.
	pstore := store.New(pdb)
	peerstore := peerstore.New(pstore)

	// Get the list of dial back peers.
	peers, err := peerstore.Peers()
	if err != nil {
		log.Error().Err(err).Msg("could not get list of dial-back peers")
		return failure
	}

	// Get the list of boot nodes addresses.
	bootNodeAddrs, err := getBootNodeAddresses(cfg.BootNodes)
	if err != nil {
		log.Error().Err(err).Msg("could not get boot node addresses")
		return failure
	}

	// Create libp2p host.
	host, err := host.New(log, cfg.Connectivity.Address, cfg.Connectivity.Port,
		host.WithPrivateKey(cfg.Connectivity.PrivateKey),
		host.WithBootNodes(bootNodeAddrs),
		host.WithDialBackPeers(peers),
		host.WithDialBackAddress(cfg.Connectivity.DialbackAddress),
		host.WithDialBackPort(cfg.Connectivity.DialbackPort),
		host.WithDialBackWebsocketPort(cfg.Connectivity.WebsocketDialbackPort),
		host.WithWebsocket(cfg.Connectivity.Websocket),
		host.WithWebsocketPort(cfg.Connectivity.WebsocketPort),
	)
	if err != nil {
		log.Error().Err(err).Str("key", cfg.Connectivity.PrivateKey).Msg("could not create host")
		return failure
	}
	defer host.Close()

	log.Info().
		Str("id", host.ID().String()).
		Strs("addresses", host.Addresses()).
		Int("boot_nodes", len(bootNodeAddrs)).
		Int("dial_back_peers", len(peers)).
		Msg("created host")

	// Set node options.
	opts := []node.Option{
		node.WithRole(role),
		node.WithConcurrency(cfg.Concurrency),
		node.WithAttributeLoading(cfg.LoadAttributes),
	}

	// If this is a worker node, initialize an executor.
	if role == blockless.WorkerNode {

		// Executor options.
		execOptions := []executor.Option{
			executor.WithWorkDir(cfg.Workspace),
			executor.WithRuntimeDir(cfg.Worker.RuntimePath),
			executor.WithExecutableName(cfg.Worker.RuntimeCLI),
		}

		if needLimiter(cfg) {
			limiter, err := limits.New(limits.WithCPUPercentage(cfg.Worker.CPUPercentageLimit), limits.WithMemoryKB(cfg.Worker.MemoryLimitKB))
			if err != nil {
				log.Error().Err(err).Msg("could not create resource limiter")
				return failure
			}

			defer func() {
				err = limiter.Shutdown()
				if err != nil {
					log.Error().Err(err).Msg("could not shutdown resource limiter")
				}
			}()

			execOptions = append(execOptions, executor.WithLimiter(limiter))
		}

		// Create an executor.
		executor, err := executor.New(log, execOptions...)
		if err != nil {
			log.Error().
				Err(err).
				Str("workspace", cfg.Workspace).
				Str("runtime_path", cfg.Worker.RuntimePath).
				Str("runtime_cli", cfg.Worker.RuntimeCLI).
				Msg("could not create an executor")
			return failure
		}

		opts = append(opts, node.WithExecutor(executor))
		opts = append(opts, node.WithWorkspace(cfg.Workspace))
	}

	// Open the pebble function database.
	fdb, err := pebble.Open(cfg.FunctionDB, &pebble.Options{Logger: &pebbleNoopLogger{}})
	if err != nil {
		log.Error().Err(err).Str("db", cfg.FunctionDB).Msg("could not open pebble function database")
		return failure
	}
	defer fdb.Close()

	functionStore := store.New(fdb)

	// Create function store.
	fstore := fstore.New(log, functionStore, cfg.Workspace)

	// If we have topics specified, use those.
	if len(cfg.Topics) > 0 {
		opts = append(opts, node.WithTopics(cfg.Topics))
	}

	// Instantiate node.
	node, err := node.New(log, host, peerstore, fstore, opts...)
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

		if cfg.Head.RestAPI == "" {
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
		server.GET("/api/v1/health", api.Health)
		server.POST("/api/v1/functions/execute", api.Execute)
		server.POST("/api/v1/functions/install", api.Install)
		server.POST("/api/v1/functions/requests/result", api.ExecutionResult)

		// Start API in a separate goroutine.
		go func() {

			log.Info().Str("port", cfg.Head.RestAPI).Msg("Node API starting")
			err := server.Start(cfg.Head.RestAPI)
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

func needLimiter(cfg *config.Config) bool {
	return cfg.Worker.CPUPercentageLimit != 1.0 || cfg.Worker.MemoryLimitKB > 0
}

func updateDirPaths(root string, cfg *config.Config) {

	workspace := cfg.Workspace
	if workspace == "" {
		workspace = filepath.Join(root, config.DefaultWorkspace)
	}
	cfg.Workspace = workspace

	peerDB := cfg.PeerDB
	if peerDB == "" {
		peerDB = filepath.Join(root, config.DefaultPeerDB)
	}
	cfg.PeerDB = peerDB

	functionDB := cfg.FunctionDB
	if functionDB == "" {
		functionDB = filepath.Join(root, config.DefaultFunctionDB)
	}
	cfg.FunctionDB = functionDB
}

func generateNodeDirName(id string) string {
	return fmt.Sprintf(".b7s_%s", id)
}
