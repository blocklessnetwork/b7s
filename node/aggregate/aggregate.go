package aggregate

import (
	"sort"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/metadata"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/models/response"
)

func Aggregate(results response.ExecutionResultMap) Results {

	total := len(results)
	if total == 0 {
		return nil
	}

	type resultStats struct {
		seen     uint
		peers    []peer.ID
		metadata map[peer.ID]metadata.Metadata
	}

	stats := make(map[execute.RuntimeOutput]resultStats)
	for executingPeer, res := range results {

		// NOTE: It might make sense to ignore stderr in comparison.
		output := res.Result.Result

		stat, ok := stats[output]
		if !ok {
			stat = resultStats{
				seen:     0,
				peers:    make([]peer.ID, 0),
				metadata: make(map[peer.ID]metadata.Metadata),
			}
		}

		stat.seen++
		stat.peers = append(stat.peers, executingPeer)
		if res.Metadata != nil {
			stat.metadata[executingPeer] = res.Metadata
		}

		stats[output] = stat
	}

	// Convert map of results to a slice.
	aggregated := make([]Result, 0, len(stats))
	for res, stat := range stats {

		aggr := Result{
			Result:    res,
			Peers:     stat.peers,
			Frequency: 100 * float64(stat.seen) / float64(total),
			Metadata:  stat.metadata,
		}

		aggregated = append(aggregated, aggr)
	}

	// Sort the slice, most frequent result first.
	sort.SliceStable(aggregated, func(i, j int) bool {
		return aggregated[i].Frequency > aggregated[j].Frequency
	})

	return aggregated
}
