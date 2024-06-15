package aggregate

import (
	"sort"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/execute"
)

type Results []Result

// Result represents the execution result along with its aggregation stats.
type Result struct {
	Result execute.RuntimeOutput `json:"result,omitempty"`
	// Peers that got this result.
	Peers []peer.ID `json:"peers,omitempty"`
	// How frequent was this result, in percentages.
	Frequency float64 `json:"frequency,omitempty"`
	// Signature of this result
	Signature []byte `json:"signature,omitempty"`
}

type resultStats struct {
	seen      uint
	peers     []peer.ID
	signature []byte
}

func Aggregate(results execute.ResultMap) Results {

	total := len(results)
	if total == 0 {
		return nil
	}

	stats := make(map[execute.RuntimeOutput]resultStats)
	for executingPeer, res := range results {

		// NOTE: It might make sense to ignore stderr in comparison.
		output := res.Result

		stat, ok := stats[output]
		if !ok {
			stats[output] = resultStats{
				seen:      0,
				peers:     make([]peer.ID, 0),
				signature: res.Signature,
			}
		}

		stat.seen++
		stat.peers = append(stat.peers, executingPeer)
		stat.signature = res.Signature

		stats[output] = stat
	}

	// Convert map of results to a slice.
	aggregated := make([]Result, 0, len(stats))
	for res, stat := range stats {

		aggr := Result{
			Result:    res,
			Peers:     stat.peers,
			Frequency: 100 * float64(stat.seen) / float64(total),
			Signature: stat.signature,
		}

		aggregated = append(aggregated, aggr)
	}

	// Sort the slice, most frequent result first.
	sort.SliceStable(aggregated, func(i, j int) bool {
		return aggregated[i].Frequency > aggregated[j].Frequency
	})

	return aggregated
}
