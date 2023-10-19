package overseer

import (
	"io"
	"time"
)

type Job struct {
	// TODO: temp, move to struct.
	Exec struct {
		Path string
		Args []string
		Env  []string
	}
	ID    string
	Stdin io.Reader
}

type JobStatus uint

const (
	StatusStarted JobStatus = iota + 1
	StatusRunning
	StatusDone
	StatusKilled
	StatusFailed
)

func (s JobStatus) String() string {
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

type JobState struct {
	Status JobStatus `json:"status,omitempty"`

	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`

	StartTime    time.Time  `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"` // TODO: Check behavior if this is not a pointer, does it get omitted?
	ObservedTime time.Time  `json:"observed_time,omitempty"`

	ExitCode *int `json:"exit_code,omitempty"`
}
