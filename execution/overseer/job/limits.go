package job

type Limits struct {
	CPUPercentage float64 `json:"cpu_percentage,omitempty"`
	MemoryLimitKB uint64  `json:"memory_limit_kb,omitempty"`
	NoExec        bool    `json:"no_exec,omitempty"`
}
