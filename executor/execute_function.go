package executor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

// ExecuteFunction will run the Blockless function defined by the execution request.
func (e *Executor) ExecuteFunction(ctx context.Context, requestID string, req execute.Request) (execute.Result, error) {

	// TODO: Check other span options and stuff.
	// TODO: More details on the execution.
	_, span := e.tracer.Start(ctx, "ExecuteFunction",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(tracing.ExecutionAttributes(requestID, req)...))
	defer span.End()

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
// execution information like standard output, standard error, exit code and resource usage.
func (e *Executor) executeFunction(requestID string, req execute.Request) (execute.RuntimeOutput, execute.Usage, error) {

	log := e.log.With().Str("request", requestID).Str("function", req.FunctionID).Logger()

	log.Info().Msg("processing execution request")

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
			log.Error().Err(err).Str("dir", paths.workdir).Msg("could not remove request working directory")
		}
	}()

	log.Debug().Str("dir", paths.workdir).Msg("working directory for the request")

	// Create command that will be executed.
	cmd := e.createCmd(paths, req)

	log.Debug().Int("env_vars_set", len(cmd.Env)).Str("cmd", cmd.String()).Msg("command ready for execution")

	out, usage, err := e.executeCommand(cmd)
	if err != nil {
		return out, execute.Usage{}, fmt.Errorf("command execution failed: %w", err)
	}

	log.Info().Msg("command executed successfully")

	return out, usage, nil
}
