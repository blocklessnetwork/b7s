package raft

import (
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
)

func (h *Handler) Execute(from peer.ID, requestID string, req execute.Request) (codes.Code, execute.Result, error) {

	h.log.Info().Msg("received an execution request")

	if !h.IsLeader() {
		_, id := h.LeaderWithID()

		h.log.Info().Str("leader", string(id)).Msg("we are not the cluster leader - dropping the request")
		return codes.NoContent, execute.Result{}, nil
	}

	h.log.Info().Msg("we are the cluster leader, executing the request")

	fsmReq := FSMLogEntry{
		RequestID: requestID,
		Origin:    from,
		Execute:   req,
	}

	payload, err := json.Marshal(fsmReq)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not serialize request for FSM: %w", err)
	}

	// Apply Raft log.
	future := h.Apply(payload, defaultRaftApplyTimeout)
	err = future.Error()
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not apply raft log: %w", err)
	}

	h.log.Info().Msg("node applied raft log")

	// Get execution result.
	response := future.Response()
	value, ok := response.(execute.Result)
	if !ok {
		fsmErr, ok := response.(error)
		if ok {
			return codes.Error, execute.Result{}, fmt.Errorf("execution encountered an error: %w", fsmErr)
		}

		return codes.Error, execute.Result{}, fmt.Errorf("unexpected FSM response format: %T", response)
	}

	h.log.Info().Msg("cluster leader executed the request")

	return codes.OK, value, nil
}
