package blockless

import (
	"runtime"
)

const (
	runtimeCLI = "bls-runtime"

	// This environment variable contains the names of environment variables
	// that are set, that originate from the execution request config.
	RuntimeEnvVarList = "BLS_LIST_VARS"
)

// RuntimeCLI returns the name of the Blockless Runtime executable.
func RuntimeCLI() string {

	cli := runtimeCLI
	if runtime.GOOS == "windows" {
		cli += ".exe"
	}

	return cli
}
