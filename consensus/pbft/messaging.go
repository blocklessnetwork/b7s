package pbft

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func (r *Replica) send(ctx context.Context, to peer.ID, msg any, protocol protocol.ID) error {

	ctx, span := r.tracer.Start(ctx, msgSendSpanName(msg, spanMessageSend))
	defer span.End()

	saveTraceContext(ctx, msg)

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	// We don't want to wait indefinitely.
	ctx, cancel := context.WithTimeout(ctx, r.cfg.NetworkTimeout)
	defer cancel()

	// Send message.
	err = r.host.SendMessageOnProtocol(ctx, to, payload, protocol)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}

// broadcast sends message to all peers in the replica set.
func (r *Replica) broadcast(ctx context.Context, msg any) error {

	ctx, span := r.tracer.Start(ctx, msgSendSpanName(msg, spanMessageBroadcast))
	defer span.End()

	saveTraceContext(ctx, msg)

	// Serialize the message.
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode record: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.NetworkTimeout)
	defer cancel()

	var (
		wg       sync.WaitGroup
		multierr *multierror.Error
		lock     sync.Mutex
	)

	for _, target := range r.peers {
		// Skip self.
		if target == r.id {
			continue
		}

		wg.Add(1)

		// Send concurrently to everyone.
		go func(peer peer.ID) {
			defer wg.Done()

			// NOTE: We could potentially retry sending if we fail once. On the other hand, somewhat unlikely they're
			// back online split second later.

			err := r.host.SendMessageOnProtocol(ctx, peer, payload, r.protocolID)
			if err != nil {

				lock.Lock()
				defer lock.Unlock()

				multierr = multierror.Append(multierr, err)
			}
		}(target)
	}

	wg.Wait()

	// If all went well, just return.
	sendErr := multierr.ErrorOrNil()
	if sendErr == nil {
		return nil
	}

	// Warn if we had more send errors than we bargained for.
	errCount := uint(multierr.Len())
	if errCount > r.f {
		r.log.Warn().Uint("f", r.f).Uint("errors", errCount).Msg("broadcast error count higher than pBFT f value")
	}

	return fmt.Errorf("could not broadcast message: %w", sendErr)
}
