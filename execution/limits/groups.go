package limits

import (
	"fmt"
	"os"
	"path"

	"github.com/containerd/cgroups/v3/cgroup2"
)

func (l *Limiter) CreateGroup(name string, opts ...LimitOption) error {

	l.Lock()
	defer l.Unlock()

	_, ok := l.limits[name]
	if ok {
		return fmt.Errorf("limit group with id %v already exists", name)
	}

	l.log.Info().Str("name", name).Msg("creating limit group")

	specs := limitsToResources(getLimits(opts...))
	cg, err := l.limits[""].manager.NewChild(name, specs)
	if err != nil {
		return fmt.Errorf("could not create cgroup (name: %v): %w", name, err)
	}

	l.limits[name] = &limitHandler{
		manager: cg,
	}

	l.log.Info().Str("name", name).Msg("limit group created")

	return nil
}

// NOTE: Non-recursive
// TODO: Check if needed at all in a mature setup.
func (l *Limiter) ListGroups() ([]string, error) {

	path := path.Join(l.mountpoint, l.cgroup)

	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open limiter root cgroup (path: %v): %w", path, err)
	}

	entries, err := dir.ReadDir(0)
	if err != nil {
		return nil, fmt.Errorf("could not list limiter root cgroup: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		names = append(names, entry.Name())
	}

	return names, nil
}

func (l *Limiter) DeleteGroup(name string) error {

	l.Lock()
	defer l.Unlock()

	// Handler exists only if this group is currently open and has a manager attached.
	// If that's not the case, we'll remove it manually.
	// This manual process may fail if the limit group has active processes assigned to it.
	lh, ok := l.limits[name]
	if !ok {

		path := path.Join(l.mountpoint, l.cgroup, name)

		l.log.Info().Str("path", path).Msg("manually deleting limit group")

		err := os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("could not remove limit group (path: %v): %w", path, err)
		}

		return nil
	}

	l.log.Info().Str("name", name).Msg("deleting limit group")

	err := lh.manager.Delete()
	if err != nil {
		return fmt.Errorf("could not delete cgroup: %w", err)
	}

	err = lh.handle.Close()
	if err != nil {
		l.log.Error().Err(err).Str("name", name).Msg("could not close file handle for limit group")
	}

	delete(l.limits, name)

	l.log.Info().Str("name", name).Msg("limit group deleted")

	return nil
}

func (l *Limiter) loadRootGroup(opts ...LimitOption) error {

	cg, err := cgroup2.Load(l.cgroup, cgroup2.WithMountpoint(l.mountpoint))
	if err != nil {
		return fmt.Errorf("could not load root cgroup: %w", err)
	}

	limits := getLimits(opts...)
	specs := limitsToResources(limits)

	err = cg.Update(specs)
	if err != nil {
		return fmt.Errorf("could not set limits for root cgroup: %w", err)
	}

	// TODO: Also open it to have a handle too. We can have both, no?
	l.limits[""] = &limitHandler{
		manager: cg,
	}

	return nil
}
