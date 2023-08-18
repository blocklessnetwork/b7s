package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/response"
)

func (n *Node) processExecute(ctx context.Context, from peer.ID, payload []byte) error {
	// We execute functions differently depending on the node role.
	if n.isHead() {
		return n.headProcessExecute(ctx, from, payload)
	}
	return n.workerProcessExecute(ctx, from, payload)
}

func (n *Node) processExecuteResponse(ctx context.Context, from peer.ID, payload []byte) error {

	// Unpack the message.
	var res response.Execute
	err := json.Unmarshal(payload, &res)
	if err != nil {
		return fmt.Errorf("could not not unpack execute response: %w", err)
	}
	res.From = from

	n.log.Debug().Str("request", res.RequestID).Str("from", from.String()).Msg("received execution response")

	key := executionResultKey(res.RequestID, from)
	n.executeResponses.Set(key, res)

	return nil
}

func executionResultKey(requestID string, peer peer.ID) string {
	return requestID + "/" + peer.String()
}

// determineOverallCode will return the resulting code from a set of results. Rules are:
// - if there's a single result, we use that results code
// - return OK if at least one result was successful
// - return error if none of the results were successful
func determineOverallCode(results map[string]execute.Result) codes.Code {

	if len(results) == 0 {
		return codes.NoContent
	}

	// For a single peer, just return its code.
	if len(results) == 1 {
		for peer := range results {
			return results[peer].Code
		}
	}

	// For multiple results - return OK if any of them succeeded.
	for _, res := range results {
		if res.Code == codes.OK {
			return codes.OK
		}
	}

	return codes.Error
}

// helper function to to convert a slice of multiaddrs to strings
func peerIDList(ids []peer.ID) []string {
	peerIDs := make([]string, 0, len(ids))
	for _, rp := range ids {
		peerIDs = append(peerIDs, rp.String())
	}
	return peerIDs
}
