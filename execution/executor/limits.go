package executor

import (
	"github.com/Maelkum/limits/limits"
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Limiter interface {
	LimitProcess(proc execute.ProcessID) error
}

// TODO: Nested nature of the groups - we want to have cumulative limit for all jobs, but also per-job limit.
type LimiterV2 interface {
	CreateGroup(id string, opts ...limits.LimitOption) error
	AssignProcessToGroup(id string, proc execute.ProcessID) error
	DeleteGroup(id string) error
	// TODO: Enumerate groups/limited processes.
}
