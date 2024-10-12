package job

type Status uint

const (
	StatusStarted Status = iota + 1
	StatusRunning
	StatusDone
	StatusKilled
	StatusFailed
)

func (s Status) String() string {
	switch s {
	case StatusStarted:
		return "started"
	case StatusRunning:
		return "running"
	case StatusDone:
		return "done"
	case StatusKilled:
		return "killed"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}
