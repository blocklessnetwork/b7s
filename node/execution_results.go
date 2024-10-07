package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus/pbft"
	"github.com/blocklessnetwork/b7s/models/execute"
)

// gatherExecutionResultsPBFT collects execution results from a PBFT cluster. This means f+1 identical results.
func (n *Node) gatherExecutionResultsPBFT(ctx context.Context, requestID string, peers []peer.ID) execute.ResultMap {

	exctx, exCancel := context.WithTimeout(ctx, n.cfg.ExecutionTimeout)
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

			key := executionResultKey(requestID, sender)
			res, ok := n.executeResponses.WaitFor(exctx, key)
			if !ok {
				return
			}

			n.log.Info().Str("peer", sender.String()).Str("request", requestID).Msg("accounted execution response from peer")

			er, ok := res[sender]
			if !ok {
				return
			}

			pub, err := sender.ExtractPublicKey()
			if err != nil {
				n.log.Error().Err(err).Msg("could not derive public key from peer ID")
				return
			}

			err = er.VerifySignature(pub)
			if err != nil {
				n.log.Error().Err(err).Msg("could not verify signature of an execution response")
				return
			}

			lock.Lock()
			defer lock.Unlock()

			reskey := peerResultMapKey(er)
			result, ok := results[reskey]
			if !ok {
				results[reskey] = aggregatedResult{
					result: er.Result,
					peers: []peer.ID{
						sender,
					},
					metadata: map[peer.ID]any{
						sender: er.Metadata,
					},
				}
				return
			}

			// Record which peers have this result, and their metadata.
			result.peers = append(result.peers, sender)
			result.metadata[sender] = er.Metadata

			if uint(len(result.peers)) >= count {
				n.log.Info().Str("request", requestID).Int("peers", len(peers)).Uint("matching_results", count).Msg("have enough matching results")
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
func (n *Node) gatherExecutionResults(ctx context.Context, requestID string, peers []peer.ID) execute.ResultMap {

	// We're willing to wait for a limited amount of time.
	exctx, exCancel := context.WithTimeout(ctx, n.cfg.ExecutionTimeout)
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

		go func() {
			defer wg.Done()
			key := executionResultKey(requestID, rp)
			// XXX: cache response.Execute
			res, ok := n.executeResponses.WaitFor(exctx, key)
			if !ok {
				return
			}

			n.log.Info().Str("peer", rp.String()).Msg("accounted execution response from peer")

			exres, ok := res[rp]
			if !ok {
				return
			}

			reslock.Lock()
			defer reslock.Unlock()
			results[rp] = exres
		}()
	}

	wg.Wait()

	return results
}

func singleNodeResultMap(id peer.ID, res execute.NodeResult) execute.ResultMap {
	return map[peer.ID]execute.NodeResult{
		id: res,
	}
}
