package executor

import (
	"github.com/blocklessnetwork/b7s/models/execute"
)

type Limiter interface {
	LimitProcess(proc execute.ProcessID) error
	ListProcesses() ([]int, error)
}
