package head

import (
	"context"
	"fmt"

	"github.com/armon/go-metrics"
	"github.com/google/uuid"

	"github.com/blessnetwork/b7s/info"
	"github.com/blessnetwork/b7s/models/execute"
	"github.com/blessnetwork/b7s/models/response"
	"github.com/blessnetwork/b7s/node"
	"github.com/blessnetwork/b7s/node/internal/waitmap"
)

type HeadNode struct {
	node.Core

	cfg Config

	rollCall           *rollCallQueue
	consensusResponses *waitmap.WaitMap[string, response.FormCluster]
	workOrderResponses *waitmap.WaitMap[string, execute.NodeResult]
}

func New(core node.Core, options ...Option) (*HeadNode, error) {

	// Initialize config.
	cfg := DefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	err := cfg.Valid()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	head := &HeadNode{
		Core: core,
		cfg:  cfg,

		rollCall:           newQueue(rollCallQueueBufferSize),
		consensusResponses: waitmap.New[string, response.FormCluster](0),
		workOrderResponses: waitmap.New[string, execute.NodeResult](executionResultCacheSize),
	}

	head.Metrics().SetGaugeWithLabels(node.NodeInfoMetric, 1,
		[]metrics.Label{
			{Name: "id", Value: head.ID()},
			{Name: "version", Value: info.VcsVersion()},
			{Name: "role", Value: "head"},
		})

	return head, nil
}

func (h *HeadNode) Run(ctx context.Context) error {
	return h.Core.Run(ctx, h.process)
}

func newRequestID() string {
	return uuid.New().String()
}
