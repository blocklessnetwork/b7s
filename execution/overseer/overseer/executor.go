package overseer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/blocklessnetwork/b7s/execution/overseer/job"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/rs/zerolog"
)

type executor struct {
	log      zerolog.Logger
	overseer *Overseer
}

func CreateExecutor(overseer *Overseer) blockless.Executor {

	exec := &executor{
		log:      overseer.log,
		overseer: overseer,
	}

	return exec
}

func (e *executor) ExecuteFunction(ctx context.Context, requestID string, req execute.Request) (execute.Result, error) {

	// TODO: Runtime
	var runtime string
	job := createJob(runtime, req)
	state, err := e.overseer.Run(job)
	if err != nil {
		e.log.Error().Err(err).Msg("job run failed")
		// NOTE: not returning here + preserving the execution error.
	}

	ru := state.ResourceUsage
	out := execute.Result{
		Result: execute.RuntimeOutput{
			Stdout: state.Stdout,
			Stderr: state.Stderr,
		},
		Usage: execute.Usage{
			WallClockTime: ru.WallClockTime,
			CPUUserTime:   ru.CPUUserTime,
			CPUSysTime:    ru.CPUSysTime,
			MemoryMaxKB:   ru.MemoryMaxKB,
		},
	}

	// This should always be the case in the case where we `run` since we've waited for the process.
	if state.ExitCode != nil {
		out.Result.ExitCode = *state.ExitCode
	} else {
		e.log.Warn().Str("request", requestID).Msg("exit code missing for executed process")
	}

	return out, nil
}

// createJob will translate the execution request to a job specification.
func createJob(runtime string, req execute.Request) job.Job {

	// Setup stdin of the command.
	var stdin io.Reader
	if req.Config.Stdin != nil {
		stdin = strings.NewReader(*req.Config.Stdin)
	}

	job := job.Job{
		Exec: job.Command{
			// TODO: Workdir handle.
			// WorkDir: req.Config.Runtime.Workdir,
			Path: runtime,
			Args: createArgs(req),
			Env:  createEnv(req),
		},
		Stdin: stdin,
	}

	return job
}

func createArgs(req execute.Request) []string {

	// Prepare CLI arguments.
	// Append the input argument first.
	var args []string
	args = append(args, req.Config.Runtime.Input)

	// Append the arguments for the runtime.
	runtimeFlags := runtimeFlags(req.Config.Runtime, req.Config.Permissions)
	args = append(args, runtimeFlags...)

	// Separate runtime arguments from the function arguments.
	args = append(args, "--")

	// Function arguments.
	for _, param := range req.Parameters {
		if param.Value != "" {
			args = append(args, param.Value)
		}
	}

	return args
}

func createEnv(req execute.Request) []string {

	// Setup environment.
	// First, pass through our environment variables.
	environ := os.Environ()

	// Second, set the variables set in the execution request.
	names := make([]string, 0, len(req.Config.Environment))
	for _, env := range req.Config.Environment {
		e := fmt.Sprintf("%s=%s", env.Name, env.Value)
		environ = append(environ, e)

		names = append(names, env.Name)
	}

	// Third and final - set the `BLS_LIST_VARS` variable with
	// the list of names of the variables from the execution request.
	blsList := strings.Join(names, ";")
	blsEnv := fmt.Sprintf("%s=%s", blockless.RuntimeEnvVarList, blsList)
	environ = append(environ, blsEnv)

	return environ
}

// TODO: Copy of the function in two places now.
// runtimeFlags returns flags that can be passed to the runtime, for example by `exec.Cmd`.
func runtimeFlags(cfg execute.BLSRuntimeConfig, permissions []string) []string {

	var flags []string

	// NOTE: The `Input` field is not a CLI flag but an argument, so it's not handled here.

	if cfg.Entry != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagEntry, cfg.Entry)
	}

	if cfg.ExecutionTime > 0 {
		flags = append(flags, "--"+execute.BLSRuntimeFlagExecutionTime, fmt.Sprint(cfg.ExecutionTime))
	}

	if cfg.DebugInfo {
		flags = append(flags, "--"+execute.BLSRuntimeFlagDebug)
	}

	if cfg.FSRoot != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagFSRoot, cfg.FSRoot)
	}

	if cfg.DriversRootPath != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagDrivers, cfg.DriversRootPath)
	}

	if cfg.Fuel > 0 {
		flags = append(flags, "--"+execute.BLSRuntimeFlagFuel, fmt.Sprint(cfg.Fuel))
	}

	if cfg.Memory > 0 {
		flags = append(flags, "--"+execute.BLSRuntimeFlagMemory, fmt.Sprint(cfg.Memory))
	}

	if cfg.Logger != "" {
		flags = append(flags, "--"+execute.BLSRuntimeFlagLogger, cfg.Logger)
	}

	for _, permission := range permissions {
		flags = append(flags, "--"+execute.BLSRuntimeFlagPermission, permission)
	}

	return flags
}
