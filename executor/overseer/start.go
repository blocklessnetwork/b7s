package overseer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"
	"sync"
	"time"
)

type handle struct {
	*sync.Mutex
	id     string
	source Job

	ctx    context.Context
	cancel context.CancelFunc

	stdout *bytes.Buffer
	stderr *bytes.Buffer

	start     time.Time
	lastCheck time.Time

	cmd *exec.Cmd
}

func (o *Overseer) Start(job Job) (any, error) {

	err := o.prepareJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not prepare job: %w", err)
	}

	h, err := o.startJob(job)
	if err != nil {
		return nil, fmt.Errorf("could not start job: %w", err)
	}

	o.Lock()
	defer o.Unlock()

	o.jobs[job.ID] = h

	return h, nil
}

func (o *Overseer) startJob(job Job) (*handle, error) {

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	ctx, cancel := context.WithCancel(context.Background())

	cmd := createCmd(ctx, job)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = job.Stdin

	o.log.Info().Str("cmd", cmd.String()).Msg("observer created command")

	start := time.Now()
	err := cmd.Start()
	if err != nil {
		cancel() // Cover this code path so the linter doesn't complain.
		return nil, fmt.Errorf("could not start job: %w", err)
	}

	handle := handle{
		Mutex:  &sync.Mutex{},
		id:     job.ID,
		source: job,
		ctx:    ctx,
		cancel: cancel,

		stdout: &stdout,
		stderr: &stderr,

		start:     start,
		lastCheck: start,

		cmd: cmd,
	}

	return &handle, nil
}

func (o *Overseer) prepareJob(job Job) error {

	workdir := path.Join(o.cfg.Workdir, job.ID)
	err := o.cfg.FS.MkdirAll(workdir, defaultFSPermissions)
	if err != nil {
		return fmt.Errorf("could not create work directory for request: %w", err)
	}

	return nil
}

func createCmd(ctx context.Context, job Job) *exec.Cmd {

	cmd := exec.CommandContext(ctx, job.Exec.Path, job.Exec.Args...)
	cmd.Env = append(cmd.Env, job.Exec.Env...)

	return cmd
}
