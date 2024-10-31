package executor

import (
	"cmp"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/armon/go-metrics"
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

// Executor provides the capabilities to run external applications.
type Executor struct {
	log     zerolog.Logger
	cfg     Config
	tracer  *tracing.Tracer
	metrics *metrics.Metrics
}

// New creates a new Executor with the specified working directory.
func New(log zerolog.Logger, options ...Option) (*Executor, error) {

	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	if cfg.RuntimePath == "" {
		return nil, errors.New("runtime path is required")
	}

	// Convert the working directory to an absolute path too.
	workdir, err := filepath.Abs(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for workspace (path: %s): %w", cfg.WorkDir, err)
	}
	cfg.WorkDir = workdir

	// We need the absolute path for the runtime, since we'll be changing
	// the working directory on execution.
	runtime, err := filepath.Abs(cfg.RuntimePath)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path for runtime (path: %s): %w", cfg.RuntimePath, err)
	}
	cfg.RuntimePath = runtime

	dir := filepath.Dir(cfg.RuntimePath)
	cfg.DriversRootPath = filepath.Join(dir, "extensions")

	// Verify the runtime path is valid.
	cliPath := filepath.Join(cfg.RuntimePath)
	_, err = cfg.FS.Stat(cliPath)
	if err != nil {
		return nil, fmt.Errorf("invalid runtime path, cli not found (path: %s): %w", cliPath, err)
	}

	e := Executor{
		log:     log,
		cfg:     cfg,
		tracer:  tracing.NewTracer(tracerName),
		metrics: cmp.Or(cfg.Metrics, metrics.Default()),
	}

	return &e, nil
}