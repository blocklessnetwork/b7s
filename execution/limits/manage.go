package limits

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (l *Limiter) AssignToGroup(name string, pid uint64) error {

	l.Lock()
	defer l.Unlock()

	lh, ok := l.limits[name]
	if !ok {
		return errors.New("unknown group")
	}

	l.log.Info().Str("name", name).Uint64("pid", pid).Msg("assigning process to limit group")

	err := lh.manager.AddProc(pid)
	if err != nil {
		return fmt.Errorf("could not assign process to group: %w", err)
	}

	l.log.Info().Str("name", name).Uint64("pid", pid).Msg("process assigned to limit group")

	return nil
}

func (l *Limiter) GetHandle(name string) (uintptr, error) {

	l.Lock()
	defer l.Unlock()

	lh, ok := l.limits[name]
	if !ok {
		return 0, errors.New("unknown group")
	}

	if lh.handle != nil {
		return lh.handle.Fd(), nil
	}

	path := filepath.Join(l.mountpoint, l.cgroup, name)
	f, err := os.Open(path)
	if err != nil {
		return -0, fmt.Errorf("could not open limit dir for reading: %w", err)
	}

	lh.handle = f

	return lh.handle.Fd(), nil
}
