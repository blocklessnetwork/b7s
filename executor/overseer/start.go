package overseer

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type Handle struct {
	*sync.Mutex
	ID     string
	source Job

	stdout       *bytes.Buffer
	outputStream *websocket.Conn
	stderr       *bytes.Buffer
	errorStream  *websocket.Conn

	start     time.Time
	lastCheck time.Time

	cmd *exec.Cmd
}

func (o *Overseer) Start(job Job) (*Handle, error) {

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

func (o *Overseer) startJob(job Job) (*Handle, error) {

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := createCmd(job)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = job.Stdin

	o.log.Info().Str("cmd", cmd.String()).Msg("observer created command")

	handle := Handle{
		Mutex:  &sync.Mutex{},
		ID:     job.ID,
		source: job,

		stdout: &stdout,
		stderr: &stderr,

		cmd: cmd,
	}

	// Create an output stream if needed.
	if job.OutputStream != "" {
		// Continue even if stdout stream cannot be established.
		outputStream, err := wsConnect(job.OutputStream)
		if err != nil {
			o.log.Error().Err(err).Str("job", job.ID).Msg("could not establish output stream")
		} else {

			ws := wsWriter{
				conn: outputStream,
				log:  o.log.With().Str("job", job.ID).Logger(),
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
			o.log.Error().Err(err).Str("job", job.ID).Msg("could not establish error stream")
		} else {

			ws := wsWriter{
				conn: errorStream,
				log:  o.log.With().Str("job", job.ID).Logger(),
			}
			handle.errorStream = errorStream

			cmd.Stderr = io.MultiWriter(&stderr, &ws)
		}
	}

	start := time.Now()
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start job: %w", err)
	}

	handle.start = start
	handle.lastCheck = start

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

func createCmd(job Job) *exec.Cmd {

	cmd := exec.Command(job.Exec.Path, job.Exec.Args...)
	cmd.Env = append(cmd.Env, job.Exec.Env...)

	return cmd
}
