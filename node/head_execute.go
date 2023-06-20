package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/request"
	"github.com/blocklessnetworking/b7s/models/response"
)

// headExecute is called on the head node. The head node will publish a roll call and delegate an execution request to chosen nodes.
// The returned map contains execution results, mapped to the peer IDs of peers who reported them.
// TODO: (raft) - return info which node was it that executed the request.
func (n *Node) headExecute(ctx context.Context, requestID string, req execute.Request) (codes.Code, execute.Result, error) {

	// TODO: (raft) if no cluster/consensus is required - request direct execution.
	quorum := 1
	if req.Config.NodeCount > 1 {
		quorum = req.Config.NodeCount
	}

	n.log.Info().Str("request_id", requestID).Int("quorum", quorum).Msg("processing execution request")

	// Phase 1. - Issue roll call to nodes.

	// Create the queue to record roll call respones.
	n.rollCall.create(requestID)
	defer n.rollCall.remove(requestID)

	err := n.issueRollCall(ctx, requestID, req.FunctionID)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not issue roll call: %w", err)
	}

	n.log.Info().Str("function_id", req.FunctionID).Str("request_id", requestID).Msg("roll call published")

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

			n.log.Warn().Str("function_id", req.FunctionID).Str("request_id", requestID).Msg("roll call timed out")
			return codes.Timeout, execute.Result{}, blockless.ErrRollCallTimeout

		case reply := <-n.rollCall.responses(requestID):

			// Check if this is the reply we want - shouldn't really happen.
			if reply.FunctionID != req.FunctionID {
				n.log.Info().Str("peer", reply.From.String()).Str("request_id", requestID).Str("function_got", reply.FunctionID).Str("function_want", req.FunctionID).Msg("skipping inadequate roll call response - wrong function")
				continue
			}

			// Check if we are connected to this peer.
			// Since we receive responses to roll call via direct messages - should not happen.
			if !n.haveConnection(reply.From) {
				n.log.Info().Str("peer", reply.From.String()).Str("request_id", reply.RequestID).Msg("skipping roll call response from unconnected peer")
				continue
			}

			n.log.Info().Str("request_id", requestID).Str("peer", reply.From.String()).Int("want_peers", quorum).Msg("roll called peer chosen for execution")

			reportingPeers = append(reportingPeers, reply.From)
			if len(reportingPeers) >= quorum {
				n.log.Info().Str("request_id", requestID).Int("want", quorum).Msg("enough peers reported for roll call")
				break rollCallResponseLoop
			}
		}
	}

	n.log.Info().Strs("peers", peerIDList(reportingPeers)).Str("function_id", req.FunctionID).Str("request_id", requestID).Msg("requesting cluster formation from peers who reported for roll call")

	// Phase 2. - Request cluster formation.

	// Create cluster formation request.
	reqCluster := request.FormCluster{
		Type:      blockless.MessageFormCluster,
		RequestID: requestID,
		Peers:     reportingPeers,
	}

	// Request execution from peers.
	err = n.sendToMany(ctx, reportingPeers, reqCluster)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not send cluster formation request to peers (function: %s, request: %s): %w", req.FunctionID, requestID, err)
	}

	// Wait for cluster confirmation messages.
	n.log.Debug().Int("want", quorum).Str("request_id", requestID).Msg("waiting for cluster to be formed")

	// We're willing to wait for a limited amount of time.
	clusterCtx, exCancel := context.WithTimeout(ctx, n.cfg.ExecutionTimeout)
	defer exCancel()

	// Wait for confirmations for cluster forming.
	bootstrapped := make(map[string]struct{})
	var rlock sync.Mutex
	var rw sync.WaitGroup
	rw.Add(len(reportingPeers))

	// Wait on peers asynchronously.
	for _, rp := range reportingPeers {
		rp := rp

		go func() {
			defer rw.Done()
			key := consensusResponseKey(requestID, rp)
			res, ok := n.consensusResponses.WaitFor(clusterCtx, key)
			if !ok {
				return
			}

			n.log.Debug().Str("request_id", requestID).Str("peer", rp.String()).Msg("accounted consensus response from roll called peer")

			fc := res.(response.FormCluster)
			if fc.Code != codes.OK {
				n.log.Debug().Str("request_id", requestID).Str("peer", rp.String()).Msg("peer failed to join consensus cluster")
				return
			}

			rlock.Lock()
			defer rlock.Unlock()
			bootstrapped[rp.String()] = struct{}{}
		}()
	}

	// Wait for results, whatever they may be.
	rw.Wait()

	// Bail if not all peers joined the cluster successfully.
	if len(bootstrapped) != quorum {
		return codes.NotAvailable, execute.Result{}, fmt.Errorf("some peers failed to join consensus cluster (have: %d, want: %d)", len(bootstrapped), quorum)
	}

	// Phase 3. - Request execution.

	// Send the execution request to peers in the cluster. Non-leaders will drop the request.
	reqExecute := request.Execute{
		Type:       blockless.MessageExecute,
		FunctionID: req.FunctionID,
		Method:     req.Method,
		Parameters: req.Parameters,
		Config:     req.Config,
		RequestID:  requestID,
	}
	err = n.sendToMany(ctx, reportingPeers, reqExecute)
	if err != nil {
		return codes.Error, execute.Result{}, fmt.Errorf("could not send execution request to peers (function: %s, request: %s): %w", req.FunctionID, requestID, err)
	}

	n.log.Debug().Int("want", quorum).Str("request_id", requestID).Msg("waiting for an execution response")

	// We're willing to wait for a limited amount of time.
	exCtx, exCancel := context.WithTimeout(ctx, n.cfg.ExecutionTimeout)
	defer exCancel()

	var (
		// We're waiting for a single execution result now, as only the cluster leader will return a result.
		result     response.Execute
		exlock     sync.Mutex
		haveResult bool
		wg         sync.WaitGroup
	)

	wg.Add(len(reportingPeers))

	// Wait on peers asynchronously.
	for _, rp := range reportingPeers {
		rp := rp

		go func() {
			defer wg.Done()

			key := executionResultKey(requestID, rp)
			res, ok := n.executeResponses.WaitFor(exCtx, key)
			if !ok {
				return
			}

			n.log.Debug().Str("request_id", requestID).Str("peer", rp.String()).Msg("accounted execution response from peer")

			er := res.(response.Execute)

			exlock.Lock()
			defer exlock.Unlock()

			haveResult = true
			result = er

			// Cancel goroutines waiting for other peers.
			exCancel()
		}()
	}

	wg.Wait()

	// We should receive an execution result back.
	// TODO: (raft) What should we do if we don't get a response back? We know which other peers should have the result.
	if !haveResult {
		return codes.NotAvailable, execute.Result{}, fmt.Errorf("no execution results received")
	}

	n.log.Info().Str("request_id", requestID).Msg("received execution response")

	return result.Code, result.Result, nil
}
