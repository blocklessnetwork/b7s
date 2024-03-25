package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestConfig_ParseCLIArgs(t *testing.T) {

	var (
		role               = "worker"
		concurrency        = uint(13)
		workspace          = "/tmp/workspace"
		bootNodes          = []string{"dummy-addr-1", "dummy-addr-2", "dummy-addr-3"}
		logLevel           = "info"
		address            = "127.0.0.1"
		port               = uint(9000)
		websocket          = true
		websocketPort      = uint(9010)
		runtimePath        = "/tmp/runtime"
		cpuPercentageLimit = 0.9
	)

	cmdline := fmt.Sprintf(
		"--role %v --concurrency %v --workspace %v --boot-nodes %v,%v --boot-nodes %v "+
			"--log-level %v --address %v --port %v --websocket %v --websocket-port %v "+
			"--runtime-path %v --cpu-percentage-limit %v",
		role, concurrency, workspace, bootNodes[0], bootNodes[1], bootNodes[2],
		logLevel, address, port, websocket, websocketPort,
		runtimePath, cpuPercentageLimit,
	)

	args, err := shlex.Split(cmdline)
	require.NoError(t, err)

	cfg, err := load(args)
	require.NoError(t, err)

	require.Equal(t, role, cfg.Role)
	require.Equal(t, concurrency, cfg.Concurrency)
	require.Equal(t, workspace, cfg.Workspace)
	require.Equal(t, bootNodes, cfg.BootNodes)
	require.Equal(t, logLevel, cfg.Log.Level)
	require.Equal(t, address, cfg.Connectivity.Address)
	require.Equal(t, port, cfg.Connectivity.Port)
	require.Equal(t, websocket, cfg.Connectivity.Websocket)
	require.Equal(t, websocketPort, cfg.Connectivity.WebsocketPort)
	require.Equal(t, runtimePath, cfg.Worker.RuntimePath)
	require.Equal(t, cpuPercentageLimit, cfg.Worker.CPUPercentageLimit)
}

func TestConfig_LoadConfigFile(t *testing.T) {

	var (
		role               = "worker"
		concurrency        = uint(27)
		workspace          = "/tmp/whatever/workspace"
		bootNodes          = []string{"dummy-addr-97", "dummy-addr-98", "dummy-addr-96"}
		logLevel           = "debug"
		address            = "127.0.0.1"
		port               = uint(9010)
		dialbackPort       = uint(9020)
		websocket          = false
		runtimePath        = "/tmp/foo/runtime"
		cpuPercentageLimit = 0.75

		cfgMap = map[string]any{
			"role":        role,
			"concurrency": concurrency,
			"workspace":   workspace,
			"boot-nodes":  bootNodes,
			"log": map[string]any{
				"level": logLevel,
			},
			"connectivity": map[string]any{
				"address":       address,
				"port":          port,
				"websocket":     websocket,
				"dialback-port": dialbackPort,
			},
			"worker": map[string]any{
				"runtime-path":         runtimePath,
				"cpu-percentage-limit": cpuPercentageLimit,
			},
		}
	)

	filepath := writeConfigFile(t, cfgMap)

	args := []string{"--config", filepath}
	cfg, err := load(args)
	require.NoError(t, err)

	require.Equal(t, role, cfg.Role)
	require.Equal(t, concurrency, cfg.Concurrency)
	require.Equal(t, workspace, cfg.Workspace)
	require.Equal(t, bootNodes, cfg.BootNodes)
	require.Equal(t, logLevel, cfg.Log.Level)
	require.Equal(t, address, cfg.Connectivity.Address)
	require.Equal(t, port, cfg.Connectivity.Port)
	require.Equal(t, dialbackPort, cfg.Connectivity.DialbackPort)
	require.Equal(t, websocket, cfg.Connectivity.Websocket)
	require.Equal(t, runtimePath, cfg.Worker.RuntimePath)
	require.Equal(t, cpuPercentageLimit, cfg.Worker.CPUPercentageLimit)
}

func TestConfig_CLIArgsWithConfigFile(t *testing.T) {

	var (
		role = "worker"

		// CLI only.
		runtimePathCLI = "/tmp/runtime"

		websocketFile     = true
		websocketPortFile = uint(9010)

		// CLI values overriding file values.

		concurrencyCLI  = uint(20)
		concurrencyFile = uint(10)

		workspaceCLI  = "/tmp/node/workspace"
		workspaceFile = "/tmp/workspace"

		bootNodesCLI  = []string{"dummy-addr-10"}
		bootNodesFile = []string{"dummy-addr-1", "dummy-addr-2", "dummy-addr-3"}

		logLevelCLI  = "debug"
		logLevelFile = "info"

		addressCLI  = "127.0.0.1"
		addressFile = "0.0.0.0"

		portCLI  = uint(10000)
		portFile = uint(9000)

		cpuPercentageLimitCLI  = 0.99
		cpuPercentageLimitFile = 0.90

		restAPICLI  = "127.0.0.1:8080"
		restAPIFile = "0.0.0.0:8080"

		cfgMap = map[string]any{
			"role":        role,
			"concurrency": concurrencyFile,
			"workspace":   workspaceFile,
			"boot-nodes":  bootNodesFile,
			"log": map[string]any{
				"level": logLevelFile,
			},
			"connectivity": map[string]any{
				"address":        addressFile,
				"port":           portFile,
				"websocket":      websocketFile,
				"websocket-port": websocketPortFile,
			},
			"worker": map[string]any{
				"cpu-percentage-limit": cpuPercentageLimitFile,
			},
			"head": map[string]any{
				"rest-api": restAPIFile,
			},
		}
	)

	filepath := writeConfigFile(t, cfgMap)

	cmdline := fmt.Sprintf(
		"--role %v --runtime-path %v --concurrency %v --workspace %v --boot-nodes %v --log-level %v --address %v --port %v --cpu-percentage-limit %v --rest-api %v --config %v",
		role,
		runtimePathCLI,
		concurrencyCLI,
		workspaceCLI,
		strings.Join(bootNodesCLI, ","),
		logLevelCLI,
		addressCLI,
		portCLI,
		cpuPercentageLimitCLI,
		restAPICLI,
		filepath,
	)

	args, err := shlex.Split(cmdline)
	require.NoError(t, err)

	cfg, err := load(args)
	require.NoError(t, err)

	// Resulting config should be the merge of specified values, with CLI overriding anything in the config file.
	require.Equal(t, role, cfg.Role)

	// Overrides.
	require.Equal(t, concurrencyCLI, cfg.Concurrency)
	require.Equal(t, workspaceCLI, cfg.Workspace)
	require.Equal(t, bootNodesCLI, cfg.BootNodes)
	require.Equal(t, logLevelCLI, cfg.Log.Level)
	require.Equal(t, addressCLI, cfg.Connectivity.Address)
	require.Equal(t, portCLI, cfg.Connectivity.Port)
	require.Equal(t, cpuPercentageLimitCLI, cfg.Worker.CPUPercentageLimit)

	// Set using one of the two methods.
	require.Equal(t, runtimePathCLI, cfg.Worker.RuntimePath)
	require.Equal(t, websocketFile, cfg.Connectivity.Websocket)
	require.Equal(t, websocketPortFile, cfg.Connectivity.WebsocketPort)
	require.Equal(t, restAPICLI, cfg.Head.API)
}

func writeConfigFile(t *testing.T, m map[string]any) string {
	t.Helper()

	data, err := yaml.Marshal(m)
	require.NoError(t, err)

	dir := t.TempDir()

	filepath := filepath.Join(dir, "config.yaml")
	err = os.WriteFile(filepath, data, 0666)
	require.NoError(t, err)

	return filepath
}
