package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

func (n *Node) processRollCall(ctx context.Context, from peer.ID, payload []byte) error {

	// Only workers respond to roll calls at the moment.
	if n.cfg.Role != blockless.WorkerNode {
		n.log.Debug().Msg("skipping roll call as a non-worker node")
		return nil
	}

	// Unpack the request.
	var req request.RollCall
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return fmt.Errorf("could not unpack request: %w", err)
	}
	req.From = from

	log := n.log.With().Str("request", req.RequestID).Str("origin", req.Origin.String()).Str("function", req.FunctionID).Logger()
	log.Debug().Msg("received roll call request")

	// TODO: (raft) temporary measure - at the moment we don't support multiple raft clusters on the same node at the same time.
	if req.ConsensusNeeded && len(n.clusters) > 0 {
		log.Warn().Msg("cannot respond to a roll call as we're already participating in one raft cluster")
		return nil
	}

	// Base response to return.
	res := response.RollCall{
		Type:       blockless.MessageRollCallResponse,
		FunctionID: req.FunctionID,
		RequestID:  req.RequestID,
		Code:       codes.Error, // CodeError by default, changed if everything goes well.
	}

	// Check if we have this function installed.
	installed, err := n.fstore.Installed(req.FunctionID)
	if err != nil {
		sendErr := n.send(ctx, req.Origin, res)
		if sendErr != nil {
			// Log send error but choose to return the original error.
			log.Error().Err(sendErr).Str("to", req.Origin.String()).Msg("could not send response")
		}

		return fmt.Errorf("could not check if function is installed: %w", err)
	}

	// We don't have this function - install it now.
	if !installed {

		log.Info().Msg("roll call but function not installed, installing now")

		err = n.installFunction(req.FunctionID, manifestURLFromCID(req.FunctionID))
		if err != nil {
			sendErr := n.send(ctx, req.Origin, res)
			if sendErr != nil {
				// Log send error but choose to return the original error.
				log.Error().Err(sendErr).Str("to", req.Origin.String()).Msg("could not send response")
			}
			return fmt.Errorf("could not install function: %w", err)
		}
	}

	log.Info().Str("origin", req.Origin.String()).Msg("reporting for roll call")

	// Send postive response.
	res.Code = codes.Accepted
	err = n.send(ctx, req.Origin, res)
	if err != nil {
		return fmt.Errorf("could not send response: %w", err)
	}

	return nil
}

// issueRollCall will create a roll call request for executing the given function.
// On successful issuance of the roll call request, we return the ID of the issued request.
func (n *Node) issueRollCall(ctx context.Context, requestID string, functionID string, consensusNeeded bool) error {

	// Create a roll call request.
	rollCall := request.RollCall{
		Type:            blockless.MessageRollCall,
		Origin:          n.host.ID(),
		FunctionID:      functionID,
		RequestID:       requestID,
		ConsensusNeeded: consensusNeeded,
	}

	// Publish the mssage.
	err := n.publish(ctx, rollCall)
	if err != nil {
		return fmt.Errorf("could not publish to topic: %w", err)
	}

	return nil
}
