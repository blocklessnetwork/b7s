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

	// NOTE: For compatibility with Windows we will manually append config param later because `shlex.Split` doesn't jive with Windows paths.
	cmdline := fmt.Sprintf(
		"--role %v --runtime-path %v --concurrency %v --workspace %v --boot-nodes %v --log-level %v --address %v --port %v --cpu-percentage-limit %v --rest-api %v",
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
	)

	args, err := shlex.Split(cmdline)
	require.NoError(t, err)

	args = append(args, "--config", fmt.Sprintf("%v", filepath))

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
	require.Equal(t, restAPICLI, cfg.Head.RestAPI)
}

func TestConfig_Environment(t *testing.T) {

	const (
		role        = "worker"
		concurrency = uint(45)
		bootNodes   = "a,b,c,d"
		topics      = "topic1,topic2,topic3"
		db          = "/tmp/db"

		logLevel = "trace"

		address               = "127.0.0.1"
		port                  = uint(9000)
		dialbackPort          = uint(9001)
		websocket             = true
		websocketPort         = uint(10000)
		websocketDialbackPort = uint(10001)

		runtimePath        = "/tmp/runtime"
		cpuPercentageLimit = float64(0.97)
		memoryLimit        = int64(512_000)
	)

	t.Setenv("B7S_Role", role)
	t.Setenv("B7S_Concurrency", fmt.Sprint(concurrency))
	t.Setenv("B7S_BootNodes", bootNodes)
	t.Setenv("B7S_Topics", topics)
	t.Setenv("B7S_DB", db)
	t.Setenv("B7S_Log_Level", logLevel)
	t.Setenv("B7S_Connectivity_Address", address)
	t.Setenv("B7S_Connectivity_Port", fmt.Sprint(port))
	t.Setenv("B7S_Connectivity_DialbackPort", fmt.Sprint(dialbackPort))
	t.Setenv("B7S_Connectivity_Websocket", fmt.Sprint(websocket))
	t.Setenv("B7S_Connectivity_WebsocketPort", fmt.Sprint(websocketPort))
	t.Setenv("B7S_Connectivity_WebsocketDialbackPort", fmt.Sprint(websocketDialbackPort))
	t.Setenv("B7S_Worker_RuntimePath", runtimePath)
	t.Setenv("B7S_Worker_CPUPercentageLimit", fmt.Sprint(cpuPercentageLimit))
	t.Setenv("B7S_Worker_MemoryLimit", fmt.Sprint(memoryLimit))

	cfg, err := Load()
	require.NoError(t, err)

	require.Equal(t, role, cfg.Role)
	require.Equal(t, concurrency, cfg.Concurrency)

	nodeList := strings.Split(bootNodes, ",")
	require.Equal(t, nodeList, cfg.BootNodes)

	topicList := strings.Split(topics, ",")
	require.Equal(t, topicList, cfg.Topics)

	require.Equal(t, db, cfg.DB)
	require.Equal(t, logLevel, cfg.Log.Level)
	require.Equal(t, address, cfg.Connectivity.Address)
	require.Equal(t, port, cfg.Connectivity.Port)
	require.Equal(t, dialbackPort, cfg.Connectivity.DialbackPort)
	require.Equal(t, websocket, cfg.Connectivity.Websocket)
	require.Equal(t, websocketPort, cfg.Connectivity.WebsocketPort)
	require.Equal(t, websocketDialbackPort, cfg.Connectivity.WebsocketDialbackPort)

	require.Equal(t, runtimePath, cfg.Worker.RuntimePath)
	require.Equal(t, cpuPercentageLimit, cfg.Worker.CPUPercentageLimit)
	require.Equal(t, memoryLimit, cfg.Worker.MemoryLimitKB)
}

func TestConfig_Priority(t *testing.T) {

	const (
		envWorkspace   = "/tmp/env/workspace"
		envAddress     = "1.1.1.1"
		envPort        = uint(1)
		envRuntimePath = "/tmp/env/runtime/path"
		envLogLevel    = "error"

		cfgWorkspace    = "/tmp/cfg/workspace"
		cfgAddress      = "2.2.2.2"
		cfgPort         = uint(2)
		cfgDialbackPort = uint(12)

		cliWorkspace = "/tmp/cli/workspace"
		cliAddress   = "3.3.3.3"
		cliLogLevel  = "debug"
	)

	var (
		cfgMap = map[string]any{
			"workspace": cfgWorkspace,
			"connectivity": map[string]any{
				"address":       cfgAddress,
				"port":          cfgPort,
				"dialback-port": cfgDialbackPort,
			},
		}
	)

	filepath := writeConfigFile(t, cfgMap)

	t.Setenv("B7S_Workspace", envWorkspace)
	t.Setenv("B7S_Connectivity_Address", envAddress)
	t.Setenv("B7S_Connectivity_Port", fmt.Sprint(envPort))
	t.Setenv("B7S_Worker_RuntimePath", envRuntimePath)
	t.Setenv("B7S_Log_Level", envLogLevel)

	// NOTE: For compatiblity with Windows we will manually append config param later because `shlex.Split` doesn't jive with Windows paths.
	cmdline := fmt.Sprintf(
		"--workspace %v --address %v --log-level %v",
		cliWorkspace,
		cliAddress,
		cliLogLevel,
	)

	args, err := shlex.Split(cmdline)
	require.NoError(t, err)

	args = append(args, "--config", fmt.Sprintf("%v", filepath))

	cfg, err := load(args)
	require.NoError(t, err)

	// Verify resulting config.
	//
	// 1. CLI flags override everything
	// 2. Config file overrides environment variables
	// 3. Environment variables
	//
	// Any config option set via lower priority methods persists if it's not overwritten.

	// This is set only via env.
	require.Equal(t, envRuntimePath, cfg.Worker.RuntimePath)

	// This is set in config file and not overwritten by CLI flags, so it should remain active.
	require.Equal(t, cfgPort, cfg.Connectivity.Port)
	// This is only set in config file.
	require.Equal(t, cfgDialbackPort, cfg.Connectivity.DialbackPort)

	// CLI flags rule everything.
	require.Equal(t, cliWorkspace, cfg.Workspace)
	require.Equal(t, cliAddress, cfg.Connectivity.Address)
	require.Equal(t, cliLogLevel, cfg.Log.Level)

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
