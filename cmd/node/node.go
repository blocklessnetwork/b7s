package main

import (
	"context"
	"fmt"

	"github.com/blessnetwork/b7s/config"
	"github.com/blessnetwork/b7s/executor"
	"github.com/blessnetwork/b7s/executor/limits"
	"github.com/blessnetwork/b7s/fstore"
	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/node"
	"github.com/blessnetwork/b7s/node/head"
	"github.com/blessnetwork/b7s/node/worker"
)

type Node interface {
	Run(context.Context) error
}

func createWorkerNode(core node.Core, store bls.Store, cfg *config.Config) (Node, func() error, error) {

	// Create function store.
	fstore := fstore.New(log.With().Str("component", "fstore").Logger(), store, cfg.Workspace)

	// Executor options.
	execOptions := []executor.Option{
		executor.WithWorkDir(cfg.Workspace),
		executor.WithRuntimeDir(cfg.Worker.RuntimePath),
		executor.WithExecutableName(cfg.Worker.RuntimeCLI),
	}

	shutdown := func() error {
		return nil
	}
	if needLimiter(cfg) {
		limiter, err := limits.New(limits.WithCPUPercentage(cfg.Worker.CPUPercentageLimit), limits.WithMemoryKB(cfg.Worker.MemoryLimitKB))
		if err != nil {
			return nil, shutdown, fmt.Errorf("could not create resource limiter")
		}

		shutdown = func() error {
			return limiter.Shutdown()
		}

		execOptions = append(execOptions, executor.WithLimiter(limiter))
	}

	// Create an executor.
	executor, err := executor.New(log.With().Str("component", "executor").Logger(), execOptions...)
	if err != nil {
		return nil, shutdown, fmt.Errorf("could not create an executor: %w", err)
	}

	worker, err := worker.New(core, fstore, executor,
		worker.AttributeLoading(cfg.LoadAttributes),
		worker.Workspace(cfg.Workspace),
	)
	if err != nil {
		return nil, shutdown, fmt.Errorf("could not create a worker node: %w", err)
	}

	return worker, shutdown, nil
}

func createHeadNode(core node.Core, cfg *config.Config) (Node, error) {

	head, err := head.New(core)
	if err != nil {
		return nil, fmt.Errorf("could not create a head node: %w", err)
	}

	return head, nil
}
