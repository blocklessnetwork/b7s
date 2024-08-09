package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/consensus/pbft"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/response"
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

			er := res.(response.Execute)

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

			exres, ok := er.Results[sender]
			if !ok {
				return
			}

			lock.Lock()
			defer lock.Unlock()

			// Equality means same result (output) and same timestamp.
			reskey := fmt.Sprintf("%+#v-%s", exres.Result.Result, er.PBFT.RequestTimestamp.String())
			result, ok := results[reskey]
			if !ok {
				results[reskey] = aggregatedResult{
					result: exres.Result,
					peers: []peer.ID{
						sender,
					},
					metadata: map[peer.ID]any{
						sender: exres.Metadata,
					},
				}
				return
			}

			// Record which peers have this result, and their metadata.
			result.peers = append(result.peers, sender)
			result.metadata[sender] = exres.Metadata

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
			res, ok := n.executeResponses.WaitFor(exctx, key)
			if !ok {
				return
			}

			n.log.Info().Str("peer", rp.String()).Msg("accounted execution response from peer")

			er := res.(response.Execute)

			exres, ok := er.Results[rp]
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
