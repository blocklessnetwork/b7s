package main

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/config"
	"github.com/blocklessnetwork/b7s/execution/executor"
	"github.com/blocklessnetwork/b7s/execution/limits"
	"github.com/blocklessnetwork/b7s/execution/overseer"
	"github.com/blocklessnetwork/b7s/models/blockless"
)

func createExecutor(log zerolog.Logger, cfg config.Config) (blockless.Executor, error) {

	totalExecutionLimits := haveExecutionLimits(cfg.Worker)
	perExecutionLimits := cfg.Worker.SupportPerExecutionLimits

	// TODO: Create noop limiter
	// TODO: Limiter gets shutdown via defer when done.
	var limiter *limits.Limiter
	if totalExecutionLimits || perExecutionLimits {

		var err error
		limiter, err = limits.New(
			log,
			cfg.Worker.CgroupMountpoint,
			cfg.Worker.CgroupName,
			limits.WithCPUPercentage(cfg.Worker.CPUPercentageLimit),
			limits.WithMemoryKB(cfg.Worker.MemoryLimitKB),
		)
		if err != nil {
			return nil, fmt.Errorf("could not create limiter: %w", err)
		}
	}

	// TODO: Switch over fully to the new executor/oversseer.
	// Unless explicitly instructed otherwise, use the classic executor, like before.
	if !cfg.Worker.UseEnhancedExecutor {

		opts := []executor.Option{
			executor.WithWorkDir(cfg.Workspace),
			executor.WithRuntimePath(cfg.Worker.RuntimePath),
		}

		// TODO: Switch over to the new limiter.
		// if totalExecutionLimits {
		// 	opts = append(opts, executor.WithLimiter(limiter))
		// }

		// Create an executor.
		executor, err := executor.New(log, opts...)
		if err != nil {
			return nil, fmt.Errorf("could not create executor: %w", err)
		}

		return executor, nil
	}

	// Use the 'enhanced' executor, backed by an overseer.

	opts := []overseer.Option{
		overseer.WithAllowlist(cfg.Worker.RuntimePath),
		overseer.WithWorkdir(cfg.Workspace),
	}

	if limiter != nil {
		opts = append(opts, overseer.WithLimiter(limiter))
	}

	ov, err := overseer.New(log, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not create overseer: %w", err)
	}

	executor := overseer.CreateExecutor(ov)

	return executor, nil
}

func haveExecutionLimits(cfg config.Worker) bool {
	return (cfg.CPUPercentageLimit > 0 && cfg.CPUPercentageLimit < 1.0) || cfg.MemoryLimitKB > 0
}
