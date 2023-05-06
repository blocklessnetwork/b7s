package executor

import (
	"fmt"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// ExecuteFunction will run the Blockless function defined by the execution request.
func (e *Executor) ExecuteFunction(requestID string, req execute.Request) (execute.Result, error) {

	// Execute the function.
	out, usage, err := e.executeFunction(requestID, req)
	if err != nil {

		res := execute.Result{
			Code:      codes.Error,
			RequestID: requestID,
			Result:    out,
			Usage:     usage,
		}

		return res, fmt.Errorf("function execution failed: %w", err)
	}

	res := execute.Result{
		Code:      codes.OK,
		RequestID: requestID,
		Result:    out,
		Usage:     usage,
	}

	return res, nil
}

// executeFunction handles the actual execution of the Blockless function. It returns the
// standard output of the blockless-cli that handled the execution. `Function`
// typically takes this output and uses it to create the appropriate execution response.
func (e *Executor) executeFunction(requestID string, req execute.Request) (execute.RuntimeOutput, execute.Usage, error) {

	e.log.Info().
		Str("id", req.FunctionID).
		Str("request_id", requestID).
		Msg("processing execution request")

	// Generate paths for execution request.
	paths := e.generateRequestPaths(requestID, req.FunctionID, req.Method)

	err := e.cfg.FS.MkdirAll(paths.workdir, defaultPermissions)
	if err != nil {
		return execute.RuntimeOutput{}, execute.Usage{}, fmt.Errorf("could not setup working directory for execution (dir: %s): %w", paths.workdir, err)
	}
	// Remove all temporary files after we're done.
	defer func() {
		err := e.cfg.FS.RemoveAll(paths.workdir)
		if err != nil {
			e.log.Error().Err(err).Str("dir", paths.workdir).
				Msg("could not remove request working directory")
		}
	}()

	e.log.Debug().
		Str("dir", paths.workdir).
		Str("request_id", requestID).
		Msg("working directory for the request")

	// Create command that will be executed.
	cmd := e.createCmd(paths, req)

	e.log.Debug().
		Str("request_id", requestID).
		Int("env_vars_set", len(cmd.Env)).
		Str("cmd", cmd.String()).
		Msg("command ready for execution")

	out, usage, err := e.executeCommand(cmd)
	if err != nil {
		return out, execute.Usage{}, fmt.Errorf("command execution failed: %w", err)
	}

	e.log.Info().
		Str("request_id", requestID).
		Msg("command executed successfully")

	return out, usage, nil
}
