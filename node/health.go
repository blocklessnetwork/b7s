package node

import (
	"context"
	"time"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/response"
)

const (
	// how often should we publish the health ping.
	healthInterval = 1 * time.Minute
)

// HealthPing will run a long running loop, publishing health signal until cancelled.
func (n *Node) HealthPing(ctx context.Context) {

	ticker := time.NewTicker(healthInterval)

	for {
		select {

		case <-ticker.C:

			msg := response.Health{
				Type: blockless.MessageHealthCheck,
				Code: response.CodeOK,
			}

			err := n.publish(ctx, msg)
			if err != nil {
				n.log.Warn().Err(err).Msg("could not publish health signal")
			}

			n.log.Debug().Msg("emitted health ping")

		case <-ctx.Done():
			n.log.Info().Msg("stopping health ping")
			return
		}
	}
}
