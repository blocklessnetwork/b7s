package executor

import (
	"context"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// executes a shell command to execute a wasm file
func Execute(ctx context.Context) ([]byte, error) {
	cmd := "echo \"hello world\""
	run := exec.Command("bash", "-c", cmd)

	run.Dir = "/tmp"
	out, err := run.Output()

	if err != nil {

		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to execute request")

		return nil, err
	}

	return out, nil
}
