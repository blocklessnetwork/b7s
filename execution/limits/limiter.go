package limits

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/containerd/cgroups/v3"
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/rs/zerolog"
)

type Limiter struct {
	*sync.Mutex
	log zerolog.Logger

	mountpoint string
	cgroup     string

	limits map[string]*limitHandler
}

type limitHandler struct {
	manager *cgroup2.Manager
	handle  *os.File
}

// New will create a new limiter that will use the specified cgroup.
// The limit options given will be set for this cgroup, and will be inherited by any nested subtrees.
func New(log zerolog.Logger, mountpoint string, parentCgroup string, opts ...LimitOption) (*Limiter, error) {

	// Check if the system supports cgroups v2.
	var haveV2 bool
	if cgroups.Mode() == cgroups.Unified {
		haveV2 = true
	}
	if !haveV2 {
		return nil, errors.New("cgroups v2 is not supported")
	}

	l := Limiter{
		log: log,

		mountpoint: mountpoint,
		cgroup:     parentCgroup,

		Mutex:  &sync.Mutex{},
		limits: make(map[string]*limitHandler),
	}

	l.log.Debug().Str("cgroup", l.cgroup).Msg("created limiter")

	err := l.loadRootGroup(opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot load root cgroup: %w", err)
	}

	return &l, nil
}

// Shutdown removes all but the root cgroup limit.
func (l *Limiter) Shutdown() error {

	l.Lock()
	defer l.Unlock()

	for id, lh := range l.limits {
		if id == "" {
			continue
		}

		if lh.handle != nil {
			err := lh.handle.Close()
			if err != nil {
				l.log.Error().Err(err).Str("id", id).Msg("could not close limits handle")
			}
		}

		err := lh.manager.Delete()
		if err != nil {
			l.log.Error().Err(err).Str("id", id).Msg("could not shutdown cgroup manager")
		}

		delete(l.limits, id)
	}

	return nil
}

func getLimits(opts ...LimitOption) Limits {
	limits := DefaultLimits
	for _, opt := range opts {
		opt(&limits)
	}
	return limits
}
