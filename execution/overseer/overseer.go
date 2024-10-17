package overseer

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

// Overseer is a lot like `Executor`, but with a more granular control. It can do the same thing an executor does, but also have
// more granular control, like starting, cancelling, stopping jobs, check in periodically to collect any stdout/stderr output etc.
type Overseer struct {
	log zerolog.Logger
	cfg Config

	*sync.Mutex
	jobs map[string]*handle
}

func New(log zerolog.Logger, options ...Option) (*Overseer, error) {

	cfg := defaultConfig
	for _, option := range options {
		option(&cfg)
	}

	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	overseer := Overseer{
		log:  log,
		cfg:  cfg,
		jobs: make(map[string]*handle),

		Mutex: &sync.Mutex{},
	}

	return &overseer, nil
}

// TODO: Add shutdown for overseer.
