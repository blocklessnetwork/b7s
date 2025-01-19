package head

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/consensus/pbft"
	"github.com/blessnetwork/b7s/models/execute"
)

// gatherExecutionResultsPBFT collects execution results from a PBFT cluster. This means f+1 identical results.
func (h *HeadNode) gatherExecutionResultsPBFT(ctx context.Context, requestID string, peers []peer.ID) execute.ResultMap {

	exctx, exCancel := context.WithTimeout(ctx, h.cfg.ExecutionTimeout)
	defer exCancel()

	type aggregatedResult struct {
		result   execute.Result
		peers    []peer.ID
		metadata map[peer.ID]any
	}

	var (
		count = pbft.MinClusterResults(uint(len(peers)))
		lock  sync.Mutex
		wg    sync.WaitGroup

		results                   = make(map[string]aggregatedResult)
		out     execute.ResultMap = make(map[peer.ID]execute.NodeResult)
	)

	// We use a map as a simple way to count identical results.
	// Equality means same result (process outputs) and same request timestamp.
	peerResultMapKey := func(res execute.NodeResult) string {
		return fmt.Sprintf("%+#v-%s", res.Result.Result, res.PBFT.RequestTimestamp.String())
	}

	wg.Add(len(peers))

	for _, rp := range peers {
		go func(sender peer.ID) {
			defer wg.Done()

			key := peerRequestKey(requestID, sender)
			res, ok := h.workOrderResponses.WaitFor(exctx, key)
			if !ok {
				return
			}

			h.Log().Info().Stringer("peer", sender).Str("request", requestID).Msg("accounted execution response from peer")

			pub, err := sender.ExtractPublicKey()
			if err != nil {
				h.Log().Error().Err(err).Msg("could not derive public key from peer ID")
				return
			}

			err = res.VerifySignature(pub)
			if err != nil {
				h.Log().Error().Err(err).Msg("could not verify signature of an execution response")
				return
			}

			lock.Lock()
			defer lock.Unlock()

			reskey := peerResultMapKey(res)
			result, ok := results[reskey]
			if !ok {
				results[reskey] = aggregatedResult{
					result: res.Result,
					peers: []peer.ID{
						sender,
					},
					metadata: map[peer.ID]any{
						sender: res.Metadata,
					},
				}
				return
			}

			// Record which peers have this result, and their metadata.
			result.peers = append(result.peers, sender)
			result.metadata[sender] = res.Metadata

			results[reskey] = result

			if uint(len(result.peers)) >= count {
				h.Log().Info().Str("request", requestID).Int("peers", len(peers)).Uint("matching_results", count).Msg("have enough matching results")
				exCancel()

				for _, peer := range result.peers {
					out[peer] = execute.NodeResult{
						Result:   result.result,
						Metadata: result.metadata[peer],
					}
				}
			}
		}(rp)
	}

	wg.Wait()

	return out
}

// gatherExecutionResults collects execution results from direct executions or raft clusters.
func (h *HeadNode) gatherExecutionResults(ctx context.Context, requestID string, peers []peer.ID) execute.ResultMap {

	// We're willing to wait for a limited amount of time.
	exctx, exCancel := context.WithTimeout(ctx, h.cfg.ExecutionTimeout)
	defer exCancel()

	var (
		results execute.ResultMap = make(map[peer.ID]execute.NodeResult)
		reslock sync.Mutex
		wg      sync.WaitGroup
	)

	wg.Add(len(peers))

	// Wait on peers asynchronously.
	for _, rp := range peers {
		rp := rp

		go func(peer peer.ID) {
			defer wg.Done()
			key := peerRequestKey(requestID, peer)
			res, ok := h.workOrderResponses.WaitFor(exctx, key)
			if !ok {
				return
			}

			h.Log().Info().Str("peer", peer.String()).Msg("accounted execution response from peer")

			reslock.Lock()
			defer reslock.Unlock()
			results[peer] = res
		}(rp)
	}

	wg.Wait()

	return results
}
