package limits

import "github.com/blocklessnetwork/b7s/models/execute"

// TODO: Remove all this
// TODO: Perhaps convert this to a legacy limiter as a specific type.
// See how did OpenTelemetry do that trick to convert a function to an interface.
func (l *Limiter) LimitProcess(proc execute.ProcessID) error {
	return l.AssignProcessToRootGroup(uint64(proc.PID))
}
