package executor

// noopLimiter is a dummy limiter used when processes run without any resource limitations.
type noopLimiter struct{}

func (n *noopLimiter) LimitProcess(pid int) error {
	return nil
}

func (n *noopLimiter) ListProcesses() ([]int, error) {
	return []int{}, nil
}
