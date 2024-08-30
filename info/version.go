package info

import (
	"runtime/debug"
	"strings"
)

// VcsVersion returns the version info, in the form of "<git-commit-hash>:<git-commit-timestamp>".
func VcsVersion() string {

	var (
		timestamp  string
		commitHash string
	)

	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range info.Settings {

			switch s.Key {
			case "vcs.time":
				timestamp = s.Value
			case "vcs.revision":
				commitHash = s.Value
			default:
				continue
			}
		}
	}

	return strings.Join([]string{commitHash, timestamp}, ":")
}
