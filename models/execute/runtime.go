package execute

import (
	"fmt"
)

// RuntimeConfig represents the CLI flags supported by the runtime
type RuntimeConfig struct {
	Entry         string `json:"entry,omitempty"`
	ExecutionTime uint64 `json:"run_time,omitempty"`
	DebugInfo     bool   `json:"debug_info,omitempty"`
	Fuel          uint64 `json:"limited_fuel,omitempty"`
	Memory        uint64 `json:"limited_memory,omitempty"`
	Logger        string `json:"runtime_logger,omitempty"`
	// Fields not allowed to be set in the request.
	Input  string `json:"-"`
	FSRoot string `json:"-"`
}

const (
	// Blockless Runtime flag names.
	RuntimeFlagEntry         = "entry"
	RuntimeFlagInput         = "input"
	RuntimeFlagExecutionTime = "run_time"
	RuntimeFlagDebug         = "debug_info"
	RuntimeFlagFuel          = "limited_fuel"
	RuntimeFlagMemory        = "limited_memory"
	RuntimeFlagFSRoot        = "fs_root_path"
	RuntimeFlagLogger        = "runtime_logger"
)

// RuntimeFlags returns flags that can be passed to the runtime, for example by `exec.Cmd`.
func RuntimeFlags(cfg RuntimeConfig) []string {

	var flags []string

	if cfg.Entry != "" {
		flags = append(flags, "--"+RuntimeFlagEntry, cfg.Entry)
	}

	if cfg.Input != "" {
		flags = append(flags, "--"+RuntimeFlagInput, cfg.Input)
	}

	if cfg.ExecutionTime > 0 {
		flags = append(flags, "--"+RuntimeFlagExecutionTime, fmt.Sprint(cfg.ExecutionTime))
	}

	if cfg.DebugInfo {
		flags = append(flags, "--"+RuntimeFlagDebug)
	}

	if cfg.FSRoot != "" {
		flags = append(flags, "--"+RuntimeFlagFSRoot, cfg.FSRoot)
	}

	if cfg.Fuel > 0 {
		flags = append(flags, "--"+RuntimeFlagFuel, fmt.Sprint(cfg.Fuel))
	}

	if cfg.Memory > 0 {
		flags = append(flags, "--"+RuntimeFlagMemory, fmt.Sprint(cfg.Memory))
	}

	if cfg.Logger != "" {
		flags = append(flags, "--"+RuntimeFlagLogger, cfg.Logger)
	}

	return flags
}
