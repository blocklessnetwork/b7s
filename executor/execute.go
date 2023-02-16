package executor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// execute handles the actual execution of the Blockless function. It returns the
// standard output of the blockless-cli that handled the execution. `Function`
// typically takes this output and uses it to create the appropriate execution response.
func (e *Executor) execute(requestID string, req execute.Request) (string, error) {

	e.log.Info().
		Str("id", req.FunctionID).
		Str("request_id", requestID).
		Msg("processing execution request")

	// Generate paths for execution request.
	paths := e.generateRequestPaths(requestID, req.FunctionID, req.Method)

	err := os.MkdirAll(paths.workdir, defaultPermissions)
	if err != nil {
		return "", fmt.Errorf("could not setup working directory for execution (dir: %s): %w", paths.workdir, err)
	}
	// Remove all temporary files after we're done.
	defer func() {
		err := os.RemoveAll(paths.workdir)
		if err != nil {
			e.log.Error().Err(err).Str("dir", paths.workdir).
				Msg("could not remove request working directory")
		}
	}()

	e.log.Debug().
		Str("dir", paths.workdir).
		Str("request_id", requestID).
		Msg("working directory for the request")

	// TODO: Super hackish, but ported. See why this is actually needed.
	err = e.writeFunctionManifest(req, paths)
	if err != nil {
		return "", fmt.Errorf("could not write function manifest: %w", err)
	}

	// Create command that will be executed.
	cmd := e.createCmd(paths, req)

	e.log.Debug().
		Str("request_id", requestID).
		Int("env_vars_set", len(cmd.Env)).
		Str("cmd", cmd.String()).
		Msg("command ready for execution")

	// Execute the command and collect output.
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	e.log.Info().
		Str("request_id", requestID).
		Msg("command executed successfully")

	return string(out), nil
}

// createCmd will create the command to be executed, prepare working directory, environment, standard input and all else.
func (e *Executor) createCmd(paths requestPaths, req execute.Request) *exec.Cmd {

	// Prepare command to be executed.
	exePath := filepath.Join(e.runtimedir, blocklessCli)
	cmd := exec.Command(exePath, paths.manifest)
	cmd.Dir = paths.workdir

	// Setup stdin of the command.
	var stdin io.Reader
	if req.Config.Stdin != nil {
		stdin = strings.NewReader(*req.Config.Stdin)
	}
	cmd.Stdin = stdin

	// Setup environment.
	// First, pass through our environment variables.
	cmd.Env = os.Environ()

	// Second, set the variables set in the execution request.
	names := make([]string, 0, len(req.Config.Environment))
	for _, env := range req.Config.Environment {
		e := fmt.Sprintf("%s=%s", env.Name, env.Value)
		cmd.Env = append(cmd.Env, e)

		names = append(names, env.Name)
	}

	// Third and final - set the `BLS_LIST_VARS` variable with
	// the list of names of the variables from the execution request.
	blsList := strings.Join(names, ";")
	blsEnv := fmt.Sprintf("%s=%s", blsListEnvName, blsList)
	cmd.Env = append(cmd.Env, blsEnv)

	return cmd
}
