package executor

type Limiter interface {
	LimitProcess(pid uint64) error
}
