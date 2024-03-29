package limits

// DefaultConfig describes the default process resource limits.
var DefaultConfig = Config{
	Cgroup:        DefaultCgroup,
	JobName:       DefaultJobObjectName,
	MemoryKB:      -1,
	CPUPercentage: DefaultCPUPercentage,
}

// Config represents the resource limits to set.
type Config struct {
	Cgroup  string // On Linux, Cgroup to use for limits.
	JobName string // On Windows, job object name to use for limits.

	MemoryKB      int64   // Maximum amount of memory allowed in kilobytes.
	CPUPercentage float64 // Percentage of the CPU time allowed.
}

// Option can be used to set limits.
type Option func(*Config)

// WithCgroup sets the path for the cgroup used for the jobs.
func WithCgroup(path string) Option {
	return func(cfg *Config) {
		cfg.Cgroup = path
	}
}

// WithJobObjectName sets the name for the job object to be used for the jobs.
func WithJobObjectName(name string) Option {
	return func(cfg *Config) {
		cfg.JobName = name
	}
}

// WithCPUPercentage sets the percentage of CPU time allowed.
func WithCPUPercentage(p float64) Option {
	return func(cfg *Config) {
		cfg.CPUPercentage = p
	}
}

// WithMemoryKB sets the max amount of memory allowed in kilobytes.
func WithMemoryKB(limit int64) Option {
	return func(cfg *Config) {
		cfg.MemoryKB = limit
	}
}
