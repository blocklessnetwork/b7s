package config

type Worker struct {
	RuntimePath               string `koanf:"runtime-path"                 flag:"runtime-path"`                 // Path to Blockless Runtime.
	SupportPerExecutionLimits bool   `koanf:"support-per-execution-limits" flag:"support-per-execution-limits"` // Should the executor support per execution resource limits.
	UseEnhancedExecutor       bool   `koanf:"use-enhanced-executor"        flag:"use-enhanced-executor"`        // Use enahnced executor, backed by an overseer.

	// Cumulative limits for all executions.
	CPUPercentageLimit float64 `koanf:"cpu-percentage-limit" flag:"cpu-percentage-limit"`
	MemoryLimitKB      int64   `koanf:"memory-limit"         flag:"memory-limit"`

	// Cgroup params that are used by the limiter.
	CgroupMountpoint string `koanf:"cgroup-mountpoint"`
	CgroupName       string `koanf:"cgroup-name"`
}
