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
	RuntimeFlagExecutionTime = "run-time"
	RuntimeFlagDebug         = "debug-info"
	RuntimeFlagFuel          = "limited-fuel"
	RuntimeFlagMemory        = "limited-memory"
	RuntimeFlagFSRoot        = "fs-root-path"
	RuntimeFlagLogger        = "runtime-logger"
	RuntimeFlagPermission    = "permission"
	RuntimeFlagEnv           = "env"
)
