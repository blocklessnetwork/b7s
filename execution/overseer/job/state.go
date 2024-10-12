package job

import (
	"time"
)

type State struct {
	Status Status `json:"status,omitempty"`

	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`

	StartTime    time.Time  `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	ObservedTime time.Time  `json:"observed_time,omitempty"`

	ExitCode *int `json:"exit_code,omitempty"`

	ResourceUsage ResourceUsage `json:"resource_usage,omitempty"`
}

// ResourceUsage represents the resource usage information for a particular execution.
type ResourceUsage struct {
	WallClockTime time.Duration `json:"wall_clock_time,omitempty"`
	CPUUserTime   time.Duration `json:"cpu_user_time,omitempty"`
	CPUSysTime    time.Duration `json:"cpu_sys_time,omitempty"`
	MemoryMaxKB   int64         `json:"memory_max_kb,omitempty"`
}
