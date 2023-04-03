package executor

// noopLimiter is a dummy limiter used when processes run without any resource limitations.
type noopLimiter struct{}

func (n *noopLimiter) LimitProcess(pid uint64) error {
	return nil
}
