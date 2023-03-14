package execute

import (
	"time"
)

// Result describes an execution result.
type Result struct {
	Code      string `json:"code"`
	Result    string `json:"result"`
	RequestID string `json:"request_id"`
	Usage     Usage  `json:"usage,omitempty"`
}

// Usage represents the resource usage information for a particular execution.
type Usage struct {
	WallClockTime time.Duration `json:"wall_clock_time,omitempty"`
	CPUUserTime   time.Duration `json:"cpu_user_time,omitempty"`
	CPUSysTime    time.Duration `json:"cpu_sys_time,omitempty"`
	MemoryMaxKB   int64         `json:"memory_max_kb,omitempty"`
}
