package pbft

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
)

// TODO (pbft): Add signatures to messages and signature verification.
// TODO (pbft): Resetting timers when an execution is done.
// TODO (pbft): View change advancing and backoff.
// TODO (pbft): Request timestamp - execution exactly once, prevent multiple/out of order executions.
// TODO (pbft): Reply format (view number etc).
// TODO (pbft): If we advance to a new view, we can clear all old messages for previous views.
// TODO (pbft): Perhaps instead of an empty digest for a NullRequest - we use an actual digest of such a request?

// Replica is a single PBFT node. Both Primary and Backup nodes are all replicas.
type Replica struct {
	// PBFT related data.
	pbftCore
	replicaState

	// Track inactivity period to trigger a view change.
	requestTimer *time.Timer

	// Components.
	log      zerolog.Logger
	host     *host.Host
	executor Executor

	// Cluster identity.
	id    peer.ID
	key   crypto.PrivKey
	peers []peer.ID

	// TODO (pbft): This is used for testing ATM, remove later.
	byzantine bool
}

// NewReplica creates a new PBFT replica.
func NewReplica(log zerolog.Logger, host *host.Host, executor Executor, peers []peer.ID, key crypto.PrivKey) (*Replica, error) {

	total := uint(len(peers))

	if total < MinimumReplicaCount {
		return nil, fmt.Errorf("too small cluster for a valid PBFT (have: %v, minimum: %v)", total, MinimumReplicaCount)
	}

	replica := Replica{
		pbftCore:     newPbftCore(total),
		replicaState: newState(),

		log:      log.With().Str("component", "pbft").Logger(),
		host:     host,
		executor: executor,

		id:    host.ID(),
		key:   key,
		peers: peers,

		byzantine: isByzantine(),
	}

	replica.log.Info().Strs("replicas", peerIDList(peers)).Uint("n", total).Uint("f", replica.f).Bool("byzantine", replica.byzantine).Msg("created PBFT replica")

	// Set the message handlers.

	// Handling messages on the PBFT protocol.
	replica.setPBFTMessageHandler()

	// Handling messages on the standard B7S protocol. We ONLY support client requests there.
	replica.setGeneralMessageHandler()

	return &replica, nil
}

func (r *Replica) Shutdown() error {
	r.stopRequestTimer()
	return nil
}

func (r *Replica) setPBFTMessageHandler() {

	// We want to only accept messages from replicas in our cluster.
	// Create a map so we can perform a faster lookup.
	pm := make(map[peer.ID]struct{})
	for _, peer := range r.peers {
		pm[peer] = struct{}{}
	}

	r.host.Host.SetStreamHandler(Protocol, func(stream network.Stream) {
		defer stream.Close()

		from := stream.Conn().RemotePeer()

		// On this protocol we only allow messages from other replicas in the cluster.
		_, known := pm[from]
		if !known {
			r.log.Info().Str("peer", from.String()).Msg("received message from a peer not in our cluster, discarding")
			return
		}

		buf := bufio.NewReader(stream)
		msg, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			r.log.Error().Err(err).Msg("error receiving direct message")
			return
		}

		r.log.Debug().Str("peer", from.String()).Msg("received message")

		err = r.processMessage(from, msg)
		if err != nil {
			r.log.Error().Err(err).Str("peer", from.String()).Msg("message processing failed")
		}
	})
}

func (r *Replica) processMessage(from peer.ID, payload []byte) error {

	// If we're acting as a byzantine replica, just don't do anything.
	// At this point we're not trying any elaborate sus behavior.
	if r.byzantine {
		return errors.New("we're a byzantine replica, ignoring received message")
	}

	msg, err := unpackMessage(payload)
	if err != nil {
		return fmt.Errorf("could not unpack message: %w", err)
	}

	// Access to individual segments (pre-prepares, prepares, commits etc) could be managed on an individual level,
	// but it's probably not worth it. This way we just do it request by request.
	r.sl.Lock()
	defer r.sl.Unlock()

	err = r.isMessageAllowed(msg)
	if err != nil {
		return fmt.Errorf("message not allowed: %w", err)
	}

	switch m := msg.(type) {

	case Request:
		return r.processRequest(from, m)

	case PrePrepare:
		return r.processPrePrepare(from, m)

	case Prepare:
		return r.processPrepare(from, m)

	case Commit:
		return r.processCommit(from, m)

	case ViewChange:
		return r.processViewChange(from, m)

	case NewView:
		return r.processNewView(from, m)
	}

	return fmt.Errorf("unexpected message type (from: %s): %T", from, msg)
}

func (r *Replica) setGeneralMessageHandler() {

	r.host.Host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer stream.Close()

		from := stream.Conn().RemotePeer()

		buf := bufio.NewReader(stream)
		payload, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			r.log.Error().Err(err).Msg("error receiving direct message")
			return
		}

		r.log.Debug().Str("peer", from.String()).Msg("received message")

		// If we're acting as a byzantine replica, just don't do anything.
		// At this point we're not trying any elaborate sus behavior.
		if r.byzantine {
			r.log.Info().Msg("we're a byzantine replica, ignoring received message")
			return
		}

		msg, err := unpackMessage(payload)
		if err != nil {
			r.log.Error().Err(err).Msg("could not unpack message")
			return
		}

		// On the general protocol we ONLY support client requests.
		req, ok := msg.(Request)
		if !ok {
			r.log.Error().Str("peer", from.String()).Type("type", msg).Msg("unexpected message type")
			return
		}

		r.sl.Lock()
		defer r.sl.Unlock()

		err = r.processRequest(from, req)
		if err != nil {
			r.log.Error().Err(err).Str("request", req.ID).Str("client", req.Origin.String()).Msg("could not process request")
			return
		}
	})
}

func (r *Replica) primaryReplicaID() peer.ID {
	return r.peers[r.currentPrimary()]
}

func (r *Replica) isPrimary() bool {
	return r.id == r.primaryReplicaID()
}

// helper function to to convert a slice of multiaddrs to strings.
func peerIDList(ids []peer.ID) []string {
	peerIDs := make([]string, 0, len(ids))
	for _, rp := range ids {
		peerIDs = append(peerIDs, rp.String())
	}
	return peerIDs
}

func (r *Replica) isMessageAllowed(msg interface{}) error {

	// If we're in an active view, we accept all but new-view messages.
	if r.activeView {

		switch msg.(type) {
		case NewView:
			return ErrActiveView
		default:
			return nil
		}
	}

	// We are in a view change. Only accept view-change and new-view messages.
	// PBFT also supports checkpoint messages, but we don't use those.
	switch msg.(type) {
	case ViewChange, NewView:
		return nil
	default:
		return ErrViewChange
	}
}

// cleanupState will discard old preprepares, prepares, commist and pending requests.
// Call this before updating the list of pending requests since for those we don't know
// in which view they were scheduled - we remove all of them.
func (r *Replica) cleanupState(thresholdView uint) {

	r.log.Debug().Uint("threshold_view", thresholdView).Msg("cleaning up replica state")

	// Cleanup pending requests.
	for id := range r.pending {
		delete(r.pending, id)
	}

	// Cleanup old preprepares.
	for id := range r.preprepares {
		if id.view < thresholdView {
			delete(r.preprepares, id)
		}
	}

	// Cleanup old prepares.
	for id := range r.prepares {
		if id.view < thresholdView {
			delete(r.prepares, id)
		}
	}

	// Cleanup old commits.
	for id := range r.commits {
		if id.view < thresholdView {
			delete(r.commits, id)
		}
	}
}

func isByzantine() bool {
	env := strings.ToLower(os.Getenv(EnvVarByzantine))

	switch env {
	case "y", "yes", "true", "1":
		return true
	default:
		return false
	}
}
