package overseer

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
	"github.com/blocklessnetwork/b7s/execution/overseer/limits"
)

type handle struct {
	*sync.Mutex
	source job.Job

	workdir      string
	stdout       *bytes.Buffer
	outputStream *websocket.Conn
	stderr       *bytes.Buffer
	errorStream  *websocket.Conn

	start     time.Time
	lastCheck time.Time

	cmd *exec.Cmd
}

func (o *Overseer) Start(job job.Job) (string, error) {

	id := uuid.NewString()
	err := o.prepareJob(id, job)
	if err != nil {
		return "", fmt.Errorf("could not prepare job: %w", err)
	}

	err = o.checkPrerequisites(job)
	if err != nil {
		return "", fmt.Errorf("prerequisites not met: %w", err)
	}

	h, err := o.startJob(id, job)
	if err != nil {
		return "", fmt.Errorf("could not start job: %w", err)
	}

	o.Lock()
	defer o.Unlock()

	o.jobs[id] = h

	return id, nil
}

func (o *Overseer) startJob(id string, job job.Job) (*handle, error) {

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd, err := o.createCmd(id, &job)
	if err != nil {
		return nil, fmt.Errorf("could not create command: %w", err)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = job.Stdin

	o.log.Info().Str("cmd", cmd.String()).Msg("overseer created command")

	handle := handle{
		Mutex:  &sync.Mutex{},
		source: job,

		workdir: cmd.Dir,
		stdout:  &stdout,
		stderr:  &stderr,

		cmd: cmd,
	}

	// Create an output stream if needed.
	if job.OutputStream != "" {
		// Continue even if stdout stream cannot be established.
		outputStream, err := wsConnect(job.OutputStream)
		if err != nil {
			o.log.Error().Err(err).Str("job", id).Msg("could not establish output stream")
		} else {

			ws := wsWriter{
				conn: outputStream,
				log:  o.log.With().Str("job", id).Logger(),
			}
			handle.outputStream = outputStream

			// Use both writers - both keep locally and stream data.
			// Websocket writer will never return errors as it's less important.
			cmd.Stdout = io.MultiWriter(&stdout, &ws)
		}
	}

	// Create an error stream too, if needed.
	if job.ErrorStream != "" {
		errorStream, err := wsConnect(job.ErrorStream)
		if err != nil {
			o.log.Error().Err(err).Str("job", id).Msg("could not establish error stream")
		} else {

			ws := wsWriter{
				conn: errorStream,
				log:  o.log.With().Str("job", id).Logger(),
			}
			handle.errorStream = errorStream

			cmd.Stderr = io.MultiWriter(&stderr, &ws)
		}
	}

	start := time.Now()
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start job: %w", err)
	}

	handle.start = start
	handle.lastCheck = start

	return &handle, nil
}

func (o *Overseer) prepareJob(id string, job job.Job) error {

	workdir := o.workdir(id)
	err := o.cfg.FS.MkdirAll(workdir, defaultFSPermissions)
	if err != nil {
		return fmt.Errorf("could not create work directory for request: %w", err)
	}

	return nil
}

func (o *Overseer) createCmd(id string, execJob *job.Job) (*exec.Cmd, error) {

	workdir := execJob.Exec.WorkDir

	// TODO: Config option - disallow setting workdir.
	if workdir == "" {
		workdir = o.workdir(id)
	}

	cmd := exec.Command(execJob.Exec.Path, execJob.Exec.Args...)
	cmd.Dir = workdir
	cmd.Env = append(cmd.Env, execJob.Exec.Env...)

	var jobLimits *job.Limits
	if execJob.Limits != nil {
		jobLimits = execJob.Limits
	}

	if o.cfg.NoChildren {
		if jobLimits == nil {
			jobLimits = &job.Limits{}
		}
		jobLimits.NoExec = true
	}

	// TODO: Set no exec for root group if required.

	if o.cfg.useLimiter {

		// TODO: Methodize this
		var (
			fd  uintptr
			err error
		)

		if jobLimits == nil {
			fd, err = o.cfg.Limiter.GetHandle("")
		} else {

			opts := getLimitOpts(*jobLimits)
			err := o.cfg.Limiter.CreateGroup(id, opts...)
			if err != nil {
				return nil, fmt.Errorf("could not create limit group for job: %w", err)
			}

			fd, err = o.cfg.Limiter.GetHandle(id)
		}

		if err != nil {
			return nil, fmt.Errorf("could not get limit group handle: %w", err)
		}

		// NOTE: Setting child limits - https://man7.org/linux/man-pages/man2/clone3.2.html
		// Relevant:
		//	This file descriptor can be obtained by opening a cgroup v2 directory using either the O_RDONLY or the O_PATH flag.
		procAttr := syscall.SysProcAttr{
			UseCgroupFD: true,
			CgroupFD:    int(fd),
		}
		cmd.SysProcAttr = &procAttr

		// TODO: Check SysProcAttr
		// Cloneflags   uintptr        // Flags for clone calls (Linux only)
		// Unshareflags uintptr        // Flags for unshare calls (Linux only)
	}

	execJob.Limits = jobLimits

	return cmd, nil
}

func (o *Overseer) workdir(id string) string {
	return filepath.Join(o.cfg.Workdir, id)
}

func getLimitOpts(jobLimits job.Limits) []limits.LimitOption {

	var opts []limits.LimitOption
	if jobLimits.CPUPercentage > 0 {
		opts = append(opts, limits.WithCPUPercentage(jobLimits.CPUPercentage))
	}

	if jobLimits.MemoryLimitKB > 0 {
		opts = append(opts, limits.WithMemoryKB(int64(jobLimits.MemoryLimitKB)))
	}

	if jobLimits.NoExec {
		opts = append(opts, limits.WithProcLimit(1))
	}

	return opts
}
