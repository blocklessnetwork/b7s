package node_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/executor"
	"github.com/blocklessnetworking/b7s/function"
	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/node"
	"github.com/blocklessnetworking/b7s/peerstore"
	"github.com/blocklessnetworking/b7s/store"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

const (
	loopback = "127.0.0.1"

	startupDelay = 5 * time.Second
)

type nodeScaffolding struct {
	dir    string
	db     *pebble.DB
	host   *host.Host
	logger zerolog.Logger
	node   *node.Node
}

func instantiateNode(t *testing.T, dirnamePattern string, role blockless.NodeRole) *nodeScaffolding {
	t.Helper()

	nodeDir := fmt.Sprintf("%v-%v-", dirnamePattern, role.String())

	// Bootstrap node directory.
	dir, err := os.MkdirTemp("", nodeDir)
	require.NoError(t, err)

	// Create logger.
	logger := createLogger(t, role)

	// Create head node libp2p host.
	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	// Create head node.
	db, node := createNode(t, dir, logger, host, role)

	ns := nodeScaffolding{
		dir:    dir,
		db:     db,
		host:   host,
		logger: logger,
		node:   node,
	}

	return &ns
}

func createLogger(t *testing.T, role blockless.NodeRole) zerolog.Logger {

	// TODO: Change where log files go.
	logName := fmt.Sprintf("%v-log.json", role.String())
	logFile, err := os.Create(logName)
	require.NoError(t, err)

	logger := zerolog.New(logFile)

	return logger
}

func createNode(t *testing.T, dir string, logger zerolog.Logger, host *host.Host, role blockless.NodeRole) (*pebble.DB, *node.Node) {
	t.Helper()

	var (
		dbDir   = path.Join(dir, "db")
		workdir = path.Join(dir, "workdir")
	)

	db, err := pebble.Open(dbDir, &pebble.Options{})
	require.NoError(t, err)

	var (
		store     = store.New(db)
		peerstore = peerstore.New(store)
		fstore    = function.NewHandler(logger, store, workdir)
	)

	opts := []node.Option{
		node.WithRole(role),
	}

	if role == blockless.WorkerNode {

		var (
			// TODO: Hardcoded right now, fix.
			runtimeDir = "/home/aco/.local/bin"
		)

		executor, err := executor.New(logger,
			executor.WithRuntimeDir(runtimeDir),
			executor.WithWorkDir(workdir),
		)
		require.NoError(t, err)

		opts = append(opts, node.WithExecutor(executor))
	}

	node, err := node.New(logger, host, store, peerstore, fstore, opts...)
	require.NoError(t, err)

	return db, node
}

// client is an external actor that can interact with the nodes.
type client struct {
	host *host.Host
}

func createClient(t *testing.T) *client {
	t.Helper()

	host, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	c := client{
		host: host,
	}

	return &c
}

func (c *client) sendInstallMessage(ctx context.Context, to peer.ID, manifestURL string, cid string) error {

	req := request.InstallFunction{
		Type:        blockless.MessageInstallFunction,
		ManifestURL: manifestURL,
		CID:         cid,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}

	err = c.host.SendMessage(ctx, to, payload)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

func (c *client) sendExecutionMessage(ctx context.Context, to peer.ID, cid string, method string) error {

	req := request.Execute{
		Type:       blockless.MessageExecute,
		FunctionID: cid,
		Method:     method,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}

	err = c.host.SendMessage(ctx, to, payload)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

type functionServer struct {
	srv *httptest.Server
}

// TODO: This can all be moved to a helper package as a separate thing!
func createFunctionServer(t *testing.T, manifestPath string, deploymentPath string, archivePath string) *functionServer {

	// Archive to serve.
	archive, err := os.ReadFile(archivePath)
	require.NoError(t, err)

	// Checksum of the archive we serve.
	checksum := sha256.Sum256(archive)

	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			path := req.URL.Path
			switch path {
			// Manifest request.
			case manifestPath:

				deploymentURL := url.URL{
					Scheme: "http",
					Host:   req.Host,
					Path:   deploymentPath,
				}

				// TODO: Not a real manifest. Lookup an existing one.
				manifest := blockless.FunctionManifest{
					FSRootPath:      "",
					DriversRootPath: "",
					LimitedFuel:     200_000_000,
					LimitedMemory:   120,
					Entry:           "",
					Deployment: blockless.Deployment{
						CID:      "dummy-cid",
						Checksum: fmt.Sprintf("%x", checksum),
						URI:      deploymentURL.String(),
					},
				}

				payload, err := json.Marshal(manifest)
				require.NoError(t, err)
				w.Write(payload)

			// Archive download request.
			case deploymentPath:
				w.Write(archive)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))

	fs := functionServer{
		srv: srv,
	}

	return &fs
}
