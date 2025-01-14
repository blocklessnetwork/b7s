package execute

const (
	BLSDefaultRuntimeEntryPoint = "_start"
)

// RuntimeConfig represents the CLI flags supported by the runtime
type BLSRuntimeConfig struct {
	Entry           string `json:"entry,omitempty"`
	ExecutionTime   uint64 `json:"run_time,omitempty"`
	DebugInfo       bool   `json:"debug_info,omitempty"`
	Fuel            uint64 `json:"limited_fuel,omitempty"`
	Memory          uint64 `json:"limited_memory,omitempty"`
	Logger          string `json:"runtime_logger,omitempty"`
	DriversRootPath string `json:"drivers_root_path,omitempty"`
	// Fields not allowed to be set in the request.
	Input  string `json:"-"`
	FSRoot string `json:"-"`
}

const (
	// Bless Runtime flag names.
	BLSRuntimeFlagEntry         = "entry"
	BLSRuntimeFlagExecutionTime = "run-time"
	BLSRuntimeFlagDebug         = "debug-info"
	BLSRuntimeFlagFuel          = "limited-fuel"
	BLSRuntimeFlagMemory        = "limited-memory"
	BLSRuntimeFlagFSRoot        = "fs-root-path"
	BLSRuntimeFlagLogger        = "runtime-logger"
	BLSRuntimeFlagPermission    = "permission"
	BLSRuntimeFlagEnv           = "env"
	BLSRuntimeFlagDrivers       = "drivers-root-path"
)
