package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus"
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
	if req.Consensus == consensus.Raft && n.haveRaftClusters() {
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

	// Send positive response.
	res.Code = codes.Accepted
	err = n.send(ctx, req.Origin, res)
	if err != nil {
		return fmt.Errorf("could not send response: %w", err)
	}

	return nil
}

func (n *Node) executeRollCall(ctx context.Context, requestID string, functionID string, nodeCount int, consensus consensus.Type) ([]peer.ID, error) {

	// Create a logger with relevant context.
	log := n.log.With().Str("request", requestID).Str("function", functionID).Int("node_count", nodeCount).Logger()

	log.Info().Msg("performing roll call for request")

	n.rollCall.create(requestID)
	defer n.rollCall.remove(requestID)

	err := n.publishRollCall(ctx, requestID, functionID, consensus)
	if err != nil {
		return nil, fmt.Errorf("could not publish roll call: %w", err)
	}

	log.Info().Msg("roll call published")

	// Limit for how long we wait for responses.
	tctx, exCancel := context.WithTimeout(ctx, n.cfg.RollCallTimeout)
	defer exCancel()

	// Peers that have reported on roll call.
	var reportingPeers []peer.ID
rollCallResponseLoop:
	for {
		// Wait for responses from nodes who want to work on the request.
		select {
		// Request timed out.
		case <-tctx.Done():

			log.Warn().Msg("roll call timed out")
			return nil, blockless.ErrRollCallTimeout

		case reply := <-n.rollCall.responses(requestID):

			// Check if this is the reply we want - shouldn't really happen.
			if reply.FunctionID != functionID {
				log.Info().Str("peer", reply.From.String()).Str("function_got", reply.FunctionID).Msg("skipping inadequate roll call response - wrong function")
				continue
			}

			// Check if we are connected to this peer.
			// Since we receive responses to roll call via direct messages - should not happen.
			if !n.haveConnection(reply.From) {
				n.log.Info().Str("peer", reply.From.String()).Msg("skipping roll call response from unconnected peer")
				continue
			}

			log.Info().Str("peer", reply.From.String()).Msg("roll called peer chosen for execution")

			reportingPeers = append(reportingPeers, reply.From)
			if len(reportingPeers) >= nodeCount {
				log.Info().Msg("enough peers reported for roll call")
				break rollCallResponseLoop
			}
		}
	}

	return reportingPeers, nil
}

// publishRollCall will create a roll call request for executing the given function.
// On successful issuance of the roll call request, we return the ID of the issued request.
func (n *Node) publishRollCall(ctx context.Context, requestID string, functionID string, consensus consensus.Type) error {

	// Create a roll call request.
	rollCall := request.RollCall{
		Type:       blockless.MessageRollCall,
		Origin:     n.host.ID(),
		FunctionID: functionID,
		RequestID:  requestID,
		Consensus:  consensus,
	}

	// Publish the mssage.
	err := n.publish(ctx, rollCall)
	if err != nil {
		return fmt.Errorf("could not publish to topic: %w", err)
	}

	return nil
}

// Temporary measure - we can't have multiple Raft clusters at this point. Remove when we remove this limitation.
func (n *Node) haveRaftClusters() bool {

	n.clusterLock.RLock()
	defer n.clusterLock.RUnlock()

	for _, cluster := range n.clusters {
		if cluster.Consensus() == consensus.Raft {
			return true
		}
	}

	return false
}
