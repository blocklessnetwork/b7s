package node

import (
	"context"
	"fmt"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

// workerExecute is called on the worker node to use its executor component to invoke the function.
func (n *Node) workerExecute(ctx context.Context, requestID string, req execute.Request) (codes.Code, execute.Result, error) {

	// Check if we have function in store.
	functionInstalled, err := n.fstore.Installed(req.FunctionID)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not lookup function in store: %w", err)
	}

	if !functionInstalled {
		return codes.NotFound, execute.Result{}, nil
	}

	// Determine if we should just execute this function, or are we part of the cluster.
	n.clusterLock.RLock()
	raftNode, ok := n.clusters[requestID]
	n.clusterLock.Unlock()

	// There's no cluster handle created - it means we only got a direct execution request.
	if !ok {
		res, err := n.executor.ExecuteFunction(requestID, req)
		if err != nil {
			return res.Code, res, fmt.Errorf("execution failed: %w", err)
		}

		return res.Code, res, nil
	}

	// We're a part of a cluster - for now acknowledge it and return an error.
	_ = raftNode
	n.log.Info().Str("request_id", requestID).Msg("execution request to be executed as part of a cluster")

	return codes.Error, execute.Result{}, fmt.Errorf("TBD: cluster execution not yet supported")
}
