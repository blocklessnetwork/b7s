package telemetry

import (
	"runtime/debug"
	"strings"
)

func vcsVersion() string {

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
