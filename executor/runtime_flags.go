package executor

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// runtimeFlags returns flags that can be passed to the runtime, for example by `exec.Cmd`.
func runtimeFlags(cfg execute.RuntimeConfig, permissions []string) []string {

	var flags []string

	// NOTE: The `Input` field is not a CLI flag but an argument, so it's not handled here.

	if cfg.Entry != "" {
		flags = append(flags, "--"+execute.RuntimeFlagEntry, cfg.Entry)
	}

	if cfg.ExecutionTime > 0 {
		flags = append(flags, "--"+execute.RuntimeFlagExecutionTime, fmt.Sprint(cfg.ExecutionTime))
	}

	if cfg.DebugInfo {
		flags = append(flags, "--"+execute.RuntimeFlagDebug)
	}

	if cfg.FSRoot != "" {
		flags = append(flags, "--"+execute.RuntimeFlagFSRoot, cfg.FSRoot)
	}

	if cfg.Fuel > 0 {
		flags = append(flags, "--"+execute.RuntimeFlagFuel, fmt.Sprint(cfg.Fuel))
	}

	if cfg.Memory > 0 {
		flags = append(flags, "--"+execute.RuntimeFlagMemory, fmt.Sprint(cfg.Memory))
	}

	if cfg.Logger != "" {
		flags = append(flags, "--"+execute.RuntimeFlagLogger, cfg.Logger)
	}

	for _, permission := range permissions {
		flags = append(flags, "--"+execute.RuntimeFlagPermission, permission)
	}

	return flags
}
