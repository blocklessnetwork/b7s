package overseer

import (
	"sync"

	"github.com/rs/zerolog"
)

// Overseer is a lot like `Executor`, but with a more granular control. It can do the same thing an executor does, but also have
// more granular control, like starting, cancelling, stopping jobs, check in periodically to collect any stdout/stderr output etc.
type Overseer struct {
	log zerolog.Logger
	cfg Config

	*sync.Mutex
	jobs map[string]*Handle
}

func New(log zerolog.Logger, cfg Config) (*Overseer, error) {

	overseer := Overseer{
		log:  log,
		cfg:  cfg,
		jobs: make(map[string]*Handle),

		Mutex: &sync.Mutex{},
	}

	return &overseer, nil
}