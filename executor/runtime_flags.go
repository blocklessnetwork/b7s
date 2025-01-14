package executor

import (
	"fmt"

	"github.com/blessnetwork/b7s/models/execute"
)

// runtimeFlags returns flags that can be passed to the runtime, for example by `exec.Cmd`.
func runtimeFlags(cfg execute.BLSRuntimeConfig, permissions []string) []string {

	var flags []string

	// NOTE: The `Input` field is not a CLI flag but an argument, so it's not handled here.

	if cfg.Entry != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagEntry, cfg.Entry)
	}

	if cfg.ExecutionTime > 0 {
		flags = append(flags, "--"+execute.BLSRuntimeFlagExecutionTime, fmt.Sprint(cfg.ExecutionTime))
	}

	if cfg.DebugInfo {
		flags = append(flags, "--"+execute.BLSRuntimeFlagDebug)
	}

	if cfg.FSRoot != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagFSRoot, cfg.FSRoot)
	}

	if cfg.DriversRootPath != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagDrivers, cfg.DriversRootPath)
	}

	if cfg.Fuel > 0 {
		flags = append(flags, "--"+execute.BLSRuntimeFlagFuel, fmt.Sprint(cfg.Fuel))
	}

	if cfg.Memory > 0 {
		flags = append(flags, "--"+execute.BLSRuntimeFlagMemory, fmt.Sprint(cfg.Memory))
	}

	if cfg.Logger != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagLogger, cfg.Logger)
	}

	for _, permission := range permissions {
		flags = append(flags, "--"+execute.BLSRuntimeFlagPermission, permission)
	}

	return flags
}
