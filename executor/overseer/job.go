package overseer

import (
	"io"
	"time"
)

type Job struct {
	ID           string    `json:"id,omitempty"`
	Exec         Command   `json:"exec,omitempty"`
	Stdin        io.Reader `json:"stdin,omitempty"`
	OutputStream string    `json:"output_stream,omitempty"`
	ErrorStream  string    `json:"error_stream,omitempty"`
}

type Command struct {
	Path string   `json:"path,omitempty"`
	Args []string `json:"args,omitempty"`
	Env  []string `json:"env,omitempty"`
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
