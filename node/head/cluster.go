package head

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blessnetwork/b7s/consensus"
	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/codes"
	"github.com/blessnetwork/b7s/models/request"
	"github.com/blessnetwork/b7s/models/response"
)

func (h *HeadNode) formCluster(ctx context.Context, requestID string, replicas []peer.ID, consensus consensus.Type) error {

	// Create cluster formation request.
	reqCluster := request.FormCluster{
		RequestID:      requestID,
		Peers:          replicas,
		Consensus:      consensus,
		ConnectionInfo: make([]peer.AddrInfo, 0, len(replicas)),
	}

	// Add connection info in case replicas don't already know of each other.
	for _, replica := range replicas {
		addrInfo := peer.AddrInfo{
			ID:    replica,
			Addrs: h.Host().Peerstore().Addrs(replica),
		}

		reqCluster.ConnectionInfo = append(reqCluster.ConnectionInfo, addrInfo)
	}

	// Request execution from peers.
	err := h.SendToMany(ctx, replicas, &reqCluster, true)
	if err != nil {
		return fmt.Errorf("could not send cluster formation request to peers: %w", err)
	}

	// Wait for cluster confirmation messages.
	h.Log().Debug().Str("request", requestID).Msg("waiting for cluster to be formed")

	// We're willing to wait for a limited amount of time.
	clusterCtx, exCancel := context.WithTimeout(ctx, h.cfg.ExecutionTimeout)
	defer exCancel()

	// Wait for confirmations for cluster forming.
	bootstrapped := make(map[string]struct{})
	var rlock sync.Mutex
	var rw sync.WaitGroup
	rw.Add(len(replicas))

	// Wait on peers asynchronously.
	for _, rp := range replicas {
		rp := rp

		go func() {
			defer rw.Done()
			key := consensusResponseKey(requestID, rp)
			fc, ok := h.consensusResponses.WaitFor(clusterCtx, key)
			if !ok {
				return
			}

			h.Log().Info().
				Stringer("peer", rp).
				Str("request", requestID).
				Msg("accounted consensus cluster response from roll called peer")

			if fc.Code != codes.OK {
				h.Log().Warn().
					Stringer("peer", rp).
					Msg("peer failed to join consensus cluster")
				return
			}

			rlock.Lock()
			defer rlock.Unlock()
			bootstrapped[rp.String()] = struct{}{}
		}()
	}

	// Wait for results, whatever they may be.
	rw.Wait()

	// Err if not all peers joined the cluster successfully.
	if len(bootstrapped) != len(replicas) {
		return fmt.Errorf("some peers failed to join consensus cluster (have: %d, want: %d)", len(bootstrapped), len(replicas))
	}

	return nil
}

func (h *HeadNode) disbandCluster(requestID string, replicas []peer.ID) error {

	msgDisband := request.DisbandCluster{
		RequestID: requestID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), consensusClusterSendTimeout)
	defer cancel()

	err := h.SendToMany(ctx, replicas, &msgDisband, true)
	if err != nil {
		return fmt.Errorf("could not send cluster disband request (request: %s): %w", requestID, err)
	}

	h.Log().Info().
		Err(err).
		Str("request", requestID).
		Strs("peers", bls.PeerIDsToStr(replicas)).
		Msg("sent cluster disband request")

	return nil
}

// processFormClusterResponse will record the cluster formation response.
func (h *HeadNode) processFormClusterResponse(ctx context.Context, from peer.ID, res response.FormCluster) error {

	h.Log().Debug().
		Stringer("from", from).
		Str("request", res.RequestID).
		Msg("received cluster formation response")

	key := consensusResponseKey(res.RequestID, from)
	h.consensusResponses.Set(key, res)

	return nil
}

func consensusResponseKey(requestID string, peer peer.ID) string {
	return requestID + "/" + peer.String()
}
