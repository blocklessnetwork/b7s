package execute

const (
	DefaultRuntimeEntryPoint = "_start"
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
	RuntimeFlagExecutionTime = "run_time"
	RuntimeFlagDebug         = "debug_info"
	RuntimeFlagFuel          = "limited_fuel"
	RuntimeFlagMemory        = "limited_memory"
	RuntimeFlagFSRoot        = "fs_root_path"
	RuntimeFlagLogger        = "runtime_logger"
)
