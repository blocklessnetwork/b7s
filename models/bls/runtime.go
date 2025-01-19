package bls

import (
	"runtime"
)

const (
	runtimeCLI = "bls-runtime"
)

// RuntimeCLI returns the name of the Bless Runtime executable.
func RuntimeCLI() string {

	cli := runtimeCLI
	if runtime.GOOS == "windows" {
		cli += ".exe"
	}

	return cli
}
