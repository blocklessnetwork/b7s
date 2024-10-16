package executor

import (
	"github.com/blocklessnetwork/b7s/execution/limits"
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Limiter interface {
	LimitProcess(proc execute.ProcessID) error
	ListProcesses() ([]int, error)
}

// TODO: Nested nature of the groups - we want to have cumulative limit for all jo
type LimiterV2 interface {
	CreateGroup(id string, opts ...limits.LimitOption) (uintptr, error)
	AssignToGroup(proc execute.ProcessID) error
	GetGroupHandle(id string) (uintptr, error)
	DeleteGroup(id string) error
	// TODO: Enumerate groups/limited processes.
}
