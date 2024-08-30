package pbft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func (r *Replica) send(to peer.ID, msg interface{}, protocol protocol.ID) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// We don't want to wait indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.NetworkTimeout)
	defer cancel()

	// Send message.
	err = r.host.SendMessageOnProtocol(ctx, to, payload, protocol)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

// broadcast sends message to all peers in the replica set.
func (r *Replica) broadcast(msg interface{}) error {

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.NetworkTimeout)
	defer cancel()

	var errGroup multierror.Group
	for _, target := range r.peers {
		target := target

		// Skip self.
		if target == r.id {
			continue
		}

		// Send concurrently to everyone.
		errGroup.Go(func() error {

			// NOTE: We could potentially retry sending if we fail once. On the other hand, somewhat unlikely they're
			// back online split second later.
			err := r.host.SendMessageOnProtocol(ctx, target, payload, r.protocolID)
			if err != nil {
				return fmt.Errorf("peer send error (peer: %v): %w", target.String(), err)
			}

			return nil
		})
	}

	// If all went well, just return.
	sendErr := errGroup.Wait()
	if sendErr.ErrorOrNil() == nil {
		return nil
	}

	// Warn if we had more send errors than we bargained for.
	errCount := uint(sendErr.Len())
	if errCount > r.f {
		r.log.Warn().Uint("f", r.f).Uint("errors", errCount).Msg("broadcast error count higher than pBFT f value")
	}

	return fmt.Errorf("could not broadcast message: %w", sendErr)
}
