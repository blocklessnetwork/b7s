package executor

import (
	"github.com/blocklessnetwork/b7s/models/execute"
)

// noopLimiter is a dummy limiter used when processes run without any resource limitations.
type noopLimiter struct{}

func (n *noopLimiter) LimitProcess(proc execute.ProcessID) error {
	return nil
}
