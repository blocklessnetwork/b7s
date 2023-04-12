//go:build integration
// +build integration

package node_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ps "github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/executor"
	"github.com/blocklessnetworking/b7s/fstore"
	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/node"
	"github.com/blocklessnetworking/b7s/peerstore"
	"github.com/blocklessnetworking/b7s/store"
	"github.com/blocklessnetworking/b7s/testing/helpers"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

const (
	loopback = "127.0.0.1"

	startupDelay = 5 * time.Second

	cleanupDisableEnv = "B7S_INTEG_CLEANUP_DISABLE"
	runtimeDirEnv     = "B7S_INTEG_RUNTIME_DIR"
)

type nodeScaffolding struct {
	dir     string
	db      *pebble.DB
	host    *host.Host
	logFile *os.File
	node    *node.Node
}

func instantiateNode(t *testing.T, dirnamePattern string, role blockless.NodeRole) *nodeScaffolding {
	t.Helper()

	nodeDir := fmt.Sprintf("%v-%v-", dirnamePattern, role.String())

	// Bootstrap node directory.
	dir, err := os.MkdirTemp("", nodeDir)
	require.NoError(t, err)

	// Create logger.
	logName := filepath.Join(dir, fmt.Sprintf("%v-log.json", role.String()))
	logFile, err := os.Create(logName)
	require.NoError(t, err)

	logger := zerolog.New(logFile)

	// Create head node libp2p host.
	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	// Create head node.
	db, node := createNode(t, dir, logger, host, role)

	ns := nodeScaffolding{
		dir:     dir,
		db:      db,
		logFile: logFile,
		host:    host,
		node:    node,
	}

	return &ns
}

func createNode(t *testing.T, dir string, logger zerolog.Logger, host *host.Host, role blockless.NodeRole) (*pebble.DB, *node.Node) {
	t.Helper()

	var (
		dbDir   = filepath.Join(dir, "db")
		workdir = filepath.Join(dir, "workdir")
	)

	db, err := pebble.Open(dbDir, &pebble.Options{})
	require.NoError(t, err)

	var (
		store     = store.New(db)
		peerstore = peerstore.New(store)
		fstore    = fstore.New(logger, store, workdir)
	)

	opts := []node.Option{
		node.WithRole(role),
	}

	if role == blockless.WorkerNode {

		runtimeDir := os.Getenv(runtimeDirEnv)

		executor, err := executor.New(logger,
			executor.WithRuntimeDir(runtimeDir),
			executor.WithWorkDir(workdir),
		)
		require.NoError(t, err)

		opts = append(opts, node.WithExecutor(executor))
	}

	node, err := node.New(logger, host, peerstore, fstore, opts...)
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

func createFunctionServer(t *testing.T, manifestPath string, deploymentPath string, archivePath string, cid string) *helpers.FunctionServer {

	manifest := blockless.FunctionManifest{
		Name:            "hello",
		FSRootPath:      "./",
		DriversRootPath: "",
		LimitedFuel:     200_000_000,
		LimitedMemory:   120,
		Entry:           "hello.wasm",
	}

	fs := helpers.CreateFunctionServer(t, manifestPath, manifest, deploymentPath, archivePath, cid)

	return fs
}

func hostAddNewPeer(t *testing.T, host *host.Host, newPeer *host.Host) {
	t.Helper()

	info := hostGetAddrInfo(t, newPeer)
	host.Peerstore().AddAddrs(info.ID, info.Addrs, ps.PermanentAddrTTL)
}

func hostGetAddrInfo(t *testing.T, host *host.Host) *peer.AddrInfo {
	t.Helper()

	addresses := host.Addresses()
	require.NotEmpty(t, addresses)

	addr := addresses[0]

	maddr, err := multiaddr.NewMultiaddr(addr)
	require.NoError(t, err)

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	require.NoError(t, err)

	return info
}

func getStreamPayload(t *testing.T, stream network.Stream, output any) {
	t.Helper()

	buf := bufio.NewReader(stream)
	payload, err := buf.ReadBytes('\n')
	require.ErrorIs(t, err, io.EOF)

	err = json.Unmarshal(payload, output)
	require.NoError(t, err)
}

func cleanupDisabled() bool {
	return os.Getenv(cleanupDisableEnv) == "yes"
}
