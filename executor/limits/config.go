package limits

import (
	"time"
)

// DefaultConfig describes the default process resource limits.
var DefaultConfig = Config{
	Cgroup:   DefaultCgroup,
	MemoryKB: -1,
	CPUTime:  time.Duration(-1),
}

// Config represents the resource limits to set.
type Config struct {
	Cgroup   string        // Cgroup to use for limits.
	MemoryKB int64         // Maximum amount of memory allowed in kilobytes.
	CPUTime  time.Duration // Total CPU time allowed.
}

// Option can be used to set limits.
type Option func(*Config)

// WithCgroup sets the path for the cgroup used for the jobs.
func WithCgroup(path string) Option {
	return func(cfg *Config) {
		cfg.Cgroup = path
	}
}

// WithCPULimit sets the total CPU time allowed.
func WithCPULimit(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.CPUTime = d
	}
}

// WithMemoryKB sets the max amount of memory allowed in kilobytes.
func WithMemoryKB(limit int64) Option {
	return func(cfg *Config) {
		cfg.MemoryKB = limit
	}
}
